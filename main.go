package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"gopkg.in/yaml.v2"
)

func main() {
	r := Results{
		Licenses: []License{},
	}
	err := fs.WalkDir(LicensesDir, ".", r.walker)
	if err != nil {
		log.Fatalf("Fatal error. %v", err)
	}
	idx, err := fuzzyfinder.Find(
		r.Licenses,
		r.titleFromIndex,
		r.getFuzzyOptions(),
	)
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
		file, err := fs.ReadFile(LicensesDir, path)
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

func (r *Results) titleFromIndex(i int) string {
	return r.Licenses[i].Title
}

func (r *Results) textFromIndex(i int) string {
	return r.Licenses[i].Text
}

func (r *Results) getFuzzyOptions() fuzzyfinder.Option {
	return fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
		if i >= 0 {
			return r.textFromIndex(i)
		}
		return ""
	})
}
