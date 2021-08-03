package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alexmullins/zip"

	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/vorteil/direktiv-apps/pkg/direktivapps"
)

const code = "com.zip.%s.error"

type FilePath struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Type string `json:"type"`
}

type RequestInput struct {
	Files    []FilePath `json:"files"`
	VarName  string     `json:"var"`  // variable name the zipfile will be saved as
	File     string     `json:"file"` // if u want to pass the array as a variable
	Password string     `json:"password"`
}

func main() {
	direktivapps.StartServer(ZipHandler)
}

func ZipHandler(w http.ResponseWriter, r *http.Request) {

	var ri RequestInput
	aid, err := direktivapps.Unmarshal(&ri, r)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshal"), err.Error())
		return
	}

	direktivapps.LogDouble(aid, "creating zip file")
	zFile, err := os.Create("attach.zip")
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-zip"), err.Error())
		return
	}
	defer zFile.Close()

	zWriter := zip.NewWriter(zFile)
	defer zWriter.Close()

	var files []FilePath
	direktivapps.LogDouble(aid, "checking if files are provided or file is provided")
	if len(ri.Files) == 0 {
		if ri.File != "" {
			data, err := ioutil.ReadFile(filepath.Join(r.Header.Get("Direktiv-TempDir"), ri.File))
			if err != nil {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "readfile"), err.Error())
				return
			}
			err = json.Unmarshal(data, &files)
			if err != nil {
				direktivapps.RespondWithError(w, fmt.Sprintf(code, "unmarshalfiles"), err.Error())
				return
			}
			ri.Files = files
		} else {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "nofiles"), err.Error())
			return
		}
	}

	direktivapps.LogDouble(aid, "creating temp files for zip")
	for _, fp := range ri.Files {
		err = func() error {
			f, err := os.Create(fp.Name)
			if err != nil {
				return err
			}
			defer f.Close()
			if fp.Type == "base64" {
				direktivapps.LogDouble(aid, "decoding base64 file before zipping")
				sDec, _ := b64.StdEncoding.DecodeString(fp.Data)
				_, err = io.Copy(f, bytes.NewReader(sDec))
				if err != nil {
					return err
				}

			} else {
				_, err = io.Copy(f, strings.NewReader(fp.Data))
				if err != nil {
					return err
				}
			}

			return nil
		}()
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "temp-file"), err.Error())
			return
		}

		direktivapps.LogDouble(aid, fmt.Sprintf("add %s to zip", fp.Name))
		err = addFileToZip(zWriter, fp.Name, ri.Password)
		if err != nil {
			direktivapps.RespondWithError(w, fmt.Sprintf(code, "add-file-to-zip"), err.Error())
			return
		}
	}

	direktivapps.LogDouble(aid, "close zip writer")

	err = zWriter.Close()
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "close-writer"), err.Error())
		return
	}

	zFile.Seek(0, 0)
	// f, err := os.Create(filepath.Join(r.Header.Get("Direktiv-TempDir"), "out", "instance", ri.VarName))
	// if err != nil {
	// 	direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-out-file"), err.Error())
	// 	return
	// }

	// io.Copy(f, zFile)
	data, err := ioutil.ReadAll(zFile)
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "read-zip"), err.Error())
		return
	}

	direktivapps.Respond(w, data)

}

func addFileToZip(zW *zip.Writer, filename, password string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	if password != "" {
		// encrypt step
		w, err := zW.Encrypt(filename, password)
		if err != nil {
			return err
		}

		// zW.Flush()
		_, err = io.Copy(w, fileToZip)
		return err
	} else {
		// Get the file information
		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Using FileInfoHeader() above only uses the basename of the file. If we want
		// to preserve the folder structure we can overwrite this with the full path.
		header.Name = filename

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zW.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, fileToZip)
		return err
	}
}
