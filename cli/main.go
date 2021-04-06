package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gernest/front"
)

var ignoreFiles = []string{".direktiv", ".git", ".github", ".gitignore", "Makefile", "cli", "pkg", "readme.md"}

func main() {
	m := front.NewMatter()
	m.Handle("+++", front.JSONHandler)

	dirs, err := ioutil.ReadDir("../")
	if err != nil {
		fmt.Println(err)
		return
	}

	addImages := ``

	for _, dir := range dirs {
		name := dir.Name()
		found := false
		for _, ignoreF := range ignoreFiles {
			if name == ignoreF {
				found = true
				break
			}
		}

		if !found {
			data, err := ioutil.ReadFile(filepath.Join("../", name, "README.md"))
			if err != nil {
				fmt.Printf("'%s/README.md' does not exist.", name)
				return
			}
			f, _, err := m.Parse(bytes.NewReader(data))
			if err != nil {
				fmt.Printf(`Please provide a frontmatter like the following at the top of the %s/README.md for the container.
+++
{
	"image": "vorteil/test"
	"desc": "This container does stuff"
}
+++
`, name)
				return
			}
			addImages += fmt.Sprintf("| %s | %s | %s |\n", fmt.Sprintf("[%s](https://hub.docker.com/r/vorteil/%s)", f["image"], name), f["desc"], fmt.Sprintf("[README](https://github.com/vorteil/direktiv-apps/tree/master/%s)", name))
		}
	}

	readme := fmt.Sprintf(`# Direktiv Apps
<em>Generated markdown from %s</em>

Simple Containers that run on Direktiv

## Containers

| Image | Description | How to Use |
| ----- | ----------- | ---------- |
%s
`, addImages, os.Getenv("SHA"))

	fmt.Printf("%s", readme)

	err = ioutil.WriteFile("/tmp/newreadme.md", []byte(readme), 0700)
	if err != nil {
		fmt.Println(err)
	}
}
