package template

import (
	"bytes"
	"text/template"
)

func Render(t string, j interface{}) (string, error) {

	t1 := template.New("t")
	buf := bytes.NewBufferString("")

	tp, err := t1.Parse(t)
	if err != nil {
		return "", err
	}

	err = tp.Execute(buf, j)
	return buf.String(), err

}
