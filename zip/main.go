package main

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
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
	zFile, err := os.Create("tmp.zip")
	if err != nil {
		direktivapps.RespondWithError(w, fmt.Sprintf(code, "create-zip"), err.Error())
		return
	}
	defer zFile.Close()

	zWriter := zip.NewWriter(zFile)
	defer zWriter.Close()

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

	// encrypt step
	w, err := zW.Encrypt(filename, password)
	if err != nil {
		return err
	}
	zW.Flush()
	_, err = io.Copy(w, fileToZip)
	return err
}
