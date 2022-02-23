package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/direktiv/direktiv-apps/pkg/reusable"
)

type requestInput struct {
	Files           []reusable.File `json:"files"`
	Scope           string          `json:"scope"`
	Name            string          `json:"name"`
	Password        string          `json:"password"`
	Return          bool            `json:"return"`
	CompressionType string          `json:"compression"`
	AdditionalProps []string        `json:"args"`
}

func createFile(f *reusable.File, p string) (string, error) {

	if p == "" {
		p = "/tmp"
	}

	r, err := f.AsReader()
	if err != nil {
		return "", err
	}
	file, err := os.Create(fmt.Sprintf("%s/%s", p, f.Name))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, r)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func zipHandler(w http.ResponseWriter, r *http.Request, ri *reusable.RequestInfo) {

	obj := new(requestInput)
	err := reusable.Unmarshal(obj, false, r)
	if err != nil {
		reusable.ReportError(w, reusable.UnmarshallError, err)
		return
	}

	fi := reusable.NewFileIterator(obj.Files, ri)

	files := []string{}
	for {
		f, err := fi.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}

		file, err := createFile(f, ri.Dir())
		if err != nil {
			reusable.ReportError(w, errForCode("zipfile"), err)
			return
		}

		defer os.Remove(file)
		files = append(files, file)

	}

	if len(obj.CompressionType) == 0 {
		obj.CompressionType = "zip"
	}
	ct := fmt.Sprintf("-t%s", obj.CompressionType)

	args := []string{"a", ct}

	if len(obj.Password) > 0 {
		pw := fmt.Sprintf("-p%s", obj.Password)
		args = append(args, pw)
	}

	var zipFile string
	if len(obj.Scope) > 0 && len(obj.Name) > 0 {
		zipFile = fmt.Sprintf("%s/out/%s/%s", ri.Dir(), obj.Scope, obj.Name)
	} else {

		randBytes := make([]byte, 16)
		rand.Read(randBytes)
		zipFile = filepath.Join(os.TempDir(), "zip"+hex.EncodeToString(randBytes)+"."+obj.CompressionType)
		defer os.RemoveAll(zipFile)
	}

	args = append(args, zipFile)
	args = append(args, files...)

	cmd := exec.Command("7z", args...)
	cmd.Stderr = ri.LogWriter()
	cmd.Stdout = ri.LogWriter()
	cmd.Dir = ri.Dir()

	ri.Logger().Infof("running %v", cmd)

	err = cmd.Run()
	if err != nil {
		reusable.ReportError(w, errForCode("zipfile"), err)
		return
	}

	if obj.Return {
		rf, err := os.Open(zipFile)
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
