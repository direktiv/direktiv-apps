package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/direktiv/direktiv-apps/pkg/reusable"

	"github.com/alexmullins/zip"
	zip2 "github.com/klauspost/compress/zip"
)

type requestInput struct {
	Files    []reusable.File `json:"files"`
	Scope    string          `json:"scope"`
	Name     string          `json:"name"`
	Password string          `json:"password"`
	Return   bool            `json:"return"`
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
	writer *zip2.Writer
}

func (uz *unencryptedZipper) Close() error {
	uz.writer.Flush()
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

	ri.Logger().Infof("receiving zip request")

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, false, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	var zipFileOut *os.File

	// create var file if set
	if len(obj.Scope) > 0 && len(obj.Name) > 0 {

		ri.Logger().Infof("setting out filew to %s/%s", obj.Scope, obj.Name)

		fn := fmt.Sprintf("%s/out/%s/%s", r.Header.Get("Direktiv-TempDir"), obj.Scope, obj.Name)
		zipFileOut, err = os.Create(fn)
		if err != nil {
			reusable.ReportError(w, reusable.UnmarshallError, err)
			return
		}
	} else {
		ri.Logger().Infof("creating zip file")
		zipFileOut, err = ioutil.TempFile("", "zip")
		if err != nil {
			reusable.ReportError(w, reusable.UnmarshallError, err)
			return
		}
	}
	defer zipFileOut.Close()

	// remove tmp file if its not set back as var
	defer func() {
		if len(obj.Scope) == 0 && len(obj.Name) == 0 {
			os.Remove(zipFileOut.Name())
		}
	}()

	ri.Logger().Infof("start zipping to %s", zipFileOut.Name())

	var wr zipper
	if len(obj.Password) > 0 {
		ri.Logger().Infof("using encrypted zipper")
		wr = &encryptedZipper{
			zip.NewWriter(zipFileOut),
			obj.Password,
		}
	} else {
		ri.Logger().Infof("using unencrypted zipper")
		wr = &unencryptedZipper{
			zip2.NewWriter(zipFileOut),
		}
	}
	defer wr.Close()

	fi := reusable.NewFileIterator(obj.Files, ri)
	for {
		f, err := fi.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}

		ri.Logger().Infof("processing %s ", f.Name)
		ze, err := wr.Create(f.Name)
		if err != nil {
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}

		c, err := f.AsString()
		if err != nil {
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}

		_, err = ze.Write([]byte(c))
		if err != nil {
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}
		ri.Logger().Infof("processed %s", f.Name)
	}
	wr.Close()

	// reponse
	if obj.Return {
		rf, err := os.Open(zipFileOut.Name())
		if err != nil {
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}
		b, err := io.ReadAll(rf)
		if err != nil {
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}
		s := base64.StdEncoding.EncodeToString(b)

		m := make(map[string]string)
		m["zip"] = s
		reusable.ReportResult(w, s)
	} else {
		reusable.ReportResult(w, "")
	}

}

func errForCode(errCode string) string {
	return fmt.Sprintf("com.zip.%s.error", errCode)
}

func main() {
	os.Mkdir("/tmp", 0777)
	reusable.StartServer(zipHandler, nil)
}
