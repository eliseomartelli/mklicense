package main

type License struct {
	Title  string `yaml:"title"`
	Spdxid string `yaml:"spdx-id"`
	Text   string
	How    string `yaml:"how"`
}

type Results struct {
	Licenses []License
}
