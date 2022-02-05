package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/direktiv/direktiv-apps/pkg/reusable"

	"github.com/alexmullins/zip"
	zip2 "github.com/klauspost/compress/zip"
)

type filePath struct {
	Name string `json:"name"`
	Data string `json:"data"`
	Type string `json:"type"`
}

type requestInput struct {
	Files    []filePath `json:"files"`
	Scope    string     `json:"scope"`
	Name     string     `json:"name"`
	Password string     `json:"password"`
}

type zipper interface {
	Create(name string) (io.Writer, error)
	Close() error
}

type encryptedZipper struct {
	writer   *zip.Writer
	password string
}

type unencryptedZipper struct {
	writer zip2.Writer
}

func (uz *unencryptedZipper) Close() error {
	return uz.writer.Close()
}

func (uz *unencryptedZipper) Create(name string) (io.Writer, error) {
	return uz.writer.Create(name)
}

func (uz *encryptedZipper) Create(name string) (io.Writer, error) {
	return uz.writer.Encrypt(name, uz.password)
}

func (uz *encryptedZipper) Close() error {
	uz.writer.Flush()
	return uz.writer.Close()
}

func zipHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, false, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	if len(obj.Scope) == 0 || len(obj.Name) == 0 {
		reusable.ReportError(w, errForCode("zipfile"), fmt.Errorf("name and scope are required"))
		return
	}

	finalName := fmt.Sprintf("%s/out/%s/%s", r.Header.Get("Direktiv-TempDir"), obj.Scope, obj.Name)
	ri.Logger().Infof("start zipping to %s", finalName)

	zipFile, err := os.Create(finalName)
	if err != nil {
		reusable.ReportError(w, errForCode("zipfile"), err)
		return
	}

	var wr zipper
	if len(obj.Password) > 0 {
		ri.Logger().Infof("using encrypted zipper")
		wr = &encryptedZipper{
			zip.NewWriter(zipFile),
			obj.Password,
		}
	} else {
		ri.Logger().Infof("using unencrypted zipper")
		wr = &unencryptedZipper{
			writer: *zip2.NewWriter(zipFile),
		}
	}

	defer wr.Close()

	// iterate through files and create tmp files for it
	for i := range obj.Files {
		fin := obj.Files[i]

		f, err := wr.Create(fin.Name)
		if err != nil {
			os.Remove(zipFile.Name())
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}

		if fin.Type == "base64" {
			ri.Logger().Infof("base64 content for %s", fin.Name)
			dec, err := base64.StdEncoding.DecodeString(fin.Data)
			if err != nil {
				os.Remove(zipFile.Name())
				reusable.ReportError(w, errForCode("zipfile"), err)
				return
			}
			f.Write(dec)
		} else {
			ri.Logger().Infof("plain text content for %s", fin.Name)
			f.Write([]byte(fin.Data))
			if err != nil {
				os.Remove(zipFile.Name())
				reusable.ReportError(w, errForCode("zipfile"), err)
				return
			}
		}
	}

	reusable.ReportResult(w, []byte("{}"))

}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.zip.%s.error", errCode)
}

func main() {
	os.Mkdir("/tmp", 0777)
	reusable.StartServer(zipHandler, nil)
}
