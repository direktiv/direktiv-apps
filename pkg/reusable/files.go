package reusable

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type File struct {
	Name        string `json:"name"`
	Data        string `json:"data"`
	Type        string `json:"type"`
	ContentType string `json:"contenttype"`
	Mode        string `json:"mode"`

	ri *RequestInfo
}

type FileIterator struct {
	idx   int
	files []File
	ri    *RequestInfo
}

func NewFileIterator(files []File, ri *RequestInfo) *FileIterator {
	return &FileIterator{
		files: files,
		ri:    ri,
	}
}

func (fi *FileIterator) Next() (*File, error) {

	if fi.idx >= len(fi.files) {
		return nil, io.EOF
	}

	f := &fi.files[fi.idx]
	f.ri = fi.ri

	fi.idx++

	return f, nil
}

func (fi *FileIterator) Reset() {
	fi.idx = 0
}

func (f *File) AsString(ri *RequestInfo) (string, error) {

	f.ri = ri

	switch f.Type {
	case TypeBase64:
		b, err := base64.StdEncoding.DecodeString(f.Data)
		return string(b), err
	case TypeFile:
		file, err := os.Open(f.Data)
		if err != nil {
			return "", err
		}
		b, err := io.ReadAll(file)
		return string(b), err
	case TypeVariable:
		v := strings.SplitN(f.Data, "/", 2)
		if len(v) != 2 {
			return "", fmt.Errorf("can not get var %s, needs format SCOPE/NAME", f.Name)
		}
		r, _, err := f.ri.ReadVar(v[0], v[1])
		if err != nil {
			return "", err
		}
		defer r.Close()
		b, err := io.ReadAll(r)
		return string(b), err
	case TypePlain:
		return f.Data, nil
	default:
		return "", fmt.Errorf("unknown type")
	}
}

func (f *File) AsBase64(ri *RequestInfo) (string, error) {

	f.ri = ri

	switch f.Type {
	case TypeBase64:
		return f.Data, nil
	case TypeFile:
		file, err := os.Open(f.Data)
		if err != nil {
			return "", err
		}
		b, err := io.ReadAll(file)
		b64 := base64.StdEncoding.EncodeToString(b)
		return b64, err
	case TypeVariable:
		v := strings.SplitN(f.Data, "/", 2)
		if len(v) != 2 {
			return "", fmt.Errorf("can not get var %s, needs format SCOPE/NAME", f.Name)
		}
		r, _, err := f.ri.ReadVar(v[0], v[1])
		if err != nil {
			return "", err
		}
		defer r.Close()
		b, err := io.ReadAll(r)
		b64 := base64.StdEncoding.EncodeToString(b)
		return b64, err
	case TypePlain:
		b := base64.StdEncoding.EncodeToString([]byte(f.Data))
		return b, nil
	default:
		return "", fmt.Errorf("unknown type")
	}

}

func (f *File) Size(ri *RequestInfo) (int, error) {

	f.ri = ri

	switch f.Type {
	case TypeBase64:
		b, err := base64.StdEncoding.DecodeString(f.Data)
		return len(b), err
	case TypeFile:
		file, err := os.Open(f.Data)
		if err != nil {
			return -1, err
		}
		b, err := io.ReadAll(file)
		return len(b), err
	case TypeVariable:
		v := strings.SplitN(f.Data, "/", 2)
		if len(v) != 2 {
			return -1, fmt.Errorf("can not get var %s, needs format SCOPE/NAME", f.Name)
		}
		r, _, err := f.ri.ReadVar(v[0], v[1])
		if err != nil {
			return -1, err
		}
		defer r.Close()
		b, err := io.ReadAll(r)
		return len(b), err
	case TypePlain:
		return len(f.Data), nil
	default:
		return -1, fmt.Errorf("unknown type")
	}

}

func (f *File) AsReader(ri *RequestInfo) (io.ReadCloser, error) {

	f.ri = ri

	switch f.Type {
	case TypeBase64:
		b, err := base64.StdEncoding.DecodeString(f.Data)
		return io.NopCloser(strings.NewReader(string(b))), err
	case TypeFile:
		return os.Open(f.Data)
	case TypeVariable:
		v := strings.SplitN(f.Data, "/", 2)
		if len(v) != 2 {
			return nil, fmt.Errorf("can not get var %s, needs format SCOPE/NAME", f.Name)
		}
		r, _, err := f.ri.ReadVar(v[0], v[1])

		return r, err
	case TypePlain:
		return io.NopCloser(strings.NewReader(f.Data)), nil
	default:
		return nil, fmt.Errorf("unknown type")
	}

}

func (f *File) AsFile(ri *RequestInfo, mode os.FileMode) (*os.File, error) {

	f.ri = ri

	file, err := ioutil.TempFile("", f.Name)
	if err != nil {
		return nil, err
	}

	if mode == 0 {
		mode = 0644

		// try to parse
		m, err := strconv.ParseUint(f.Mode, 8, 32)
		if err == nil {
			mode = fs.FileMode(m)
		}
	}

	err = os.Chmod(file.Name(), mode)
	if err != nil {
		return nil, err
	}

	script, err := f.AsReader(ri)
	if err != nil {
		return nil, err
	}
	defer script.Close()

	_, err = io.Copy(file, script)
	if err != nil {
		return nil, err
	}

	file.Seek(0, io.SeekStart)

	return file, nil
}
