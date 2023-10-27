package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Name string
	ContentType string
	Layout string `yaml:"layout"`
	Title string `yaml:"title"`
	Permalink string `yaml:"permalink"`
	Description string `yaml:"description"`
	Author string `yaml:"author"`
	AuthorLink string `yaml:"author_link"`
	Date string `yaml:"date"`
	Categories string `yaml:"categories"`
	Body string `yaml:"body"`
	Url string `yaml:"url"`
}

func (p *Page) SetFromFile(fpath string) error {
	f, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("Error opening file %s: %w", fpath, err)
	}
	fscan := bufio.NewScanner(f)
	fscan.Split(bufio.ScanLines)

	foundHeader := false
	gotHeader := false
	header := ""
	body := ""
	for fscan.Scan() {
		if fscan.Text() == "---" {
			if foundHeader {
				gotHeader = true
				continue
			} else {
				foundHeader = true
				continue
			}
		}
		if !gotHeader {
			header = header + fscan.Text() + "\n"
		} else {
			body = body + fscan.Text() + "\n"
		}
	}
	f.Close()

	p.Name = strings.Replace(filepath.Base(fpath), ".markdown", "", -1)
	p.Name = strings.Replace(p.Name, ".html", "", -1)

	if strings.HasSuffix(fpath, ".markdown") {
		p.ContentType = "markdown"
	} else {
		p.ContentType = "html"
	}

	p.Body = body

	if err := yaml.Unmarshal([]byte(header), p); err != nil {
		return fmt.Errorf("Error setting page from YAML: %w", err)
	}

	return nil
}

func (p *Page) Validate() error {
	return nil
}
