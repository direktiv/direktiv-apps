package template

import (
	"encoding/json"
	"testing"
)

type inJson struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

func TestRender(t *testing.T) {

	g := inJson{
		Name:  "testme",
		Items: []string{"item1", "items2", "items3"},
	}
	b, _ := json.Marshal(g)

	var g2 interface{}
	json.Unmarshal(b, &g2)

	o, err := Render("hello {{.name}}! {{range .items}}{{.}} {{end}}", g2)

	t.Logf("Output: %v, error %v", o, err)

}
