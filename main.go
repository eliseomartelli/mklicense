package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"gopkg.in/yaml.v2"
)

type License struct {
	Title  string `yaml:"title"`
	Spdxid string `yaml:"spdx-id"`
	Text   string
}

type Results struct {
	Licenses []License
}

//go:embed licenses
var licensesFiles embed.FS

func main() {
	r := Results{
		Licenses: []License{},
	}
	err := fs.WalkDir(licensesFiles, ".", r.walker)
	if err != nil {
		log.Fatalf("Fatal error. %v", err)
	}
	idx, err := fuzzyfinder.Find(
		r.Licenses,
		func(i int) string {
			return r.Licenses[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("%s", r.Licenses[i].Text)
		}))
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("LICENSE", []byte(r.Licenses[idx].Text), 0660)
	fmt.Println("Open LICENSE file and change author.")
}

func (r *Results) walker(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		// It's a file!
		file, err := fs.ReadFile(licensesFiles, path)
		if err != nil {
			return err
		}
		processed := strings.Split(string(file), "---")
		license := License{}
		err = yaml.Unmarshal([]byte(processed[1]), &license)
		if err != nil {
			return err
		}
		license.Text = strings.TrimSpace(processed[2])
		r.Licenses = append(r.Licenses, license)
	}
	return nil
}
