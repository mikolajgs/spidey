package main

import (
	"os"
	"fmt"
	"path/filepath"
	"strings"
)

type Include struct {
	Name string
	Body string
	ContentType string
}

func (i *Include) SetFromFile(fpath string) error {
	body, err := os.ReadFile(fpath)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %w", fpath, err)
	}

	i.Name = strings.Replace(filepath.Base(fpath), ".markdown", "", -1)
	i.Name = strings.Replace(i.Name, ".html", "", -1)

	if strings.HasSuffix(fpath, ".markdown") {
		i.ContentType = "markdown"
	} else {
		i.ContentType = "html"
	}

	i.Body = string(body)

	return nil
}