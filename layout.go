package main

import (
	"os"
	"fmt"
	"path/filepath"
	"strings"
)

type Layout struct {
	Name string
	Body string
	ContentType string
}

func (l *Layout) SetFromFile(fpath string) error {
	body, err := os.ReadFile(fpath)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %w", fpath, err)
	}

	l.Name = strings.Replace(filepath.Base(fpath), ".markdown", "", -1)
	l.Name = strings.Replace(l.Name, ".html", "", -1)

	if strings.HasSuffix(fpath, ".markdown") {
		l.ContentType = "markdown"
	} else {
		l.ContentType = "html"
	}

	l.Body = string(body)

	return nil
}