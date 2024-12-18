package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Website struct {
	SourcePath string
	Config     *Config

	Pages     map[string]*Page
	PageNames []string

	Layouts     map[string]*Layout
	LayoutNames []string

	Includes     map[string]*Include
	IncludeNames []string

	Posts      map[string]*Page
	PostsNames []string
}

func (w *Website) Init() error {
	if err := w.initConfig(); err != nil {
		return fmt.Errorf("Error initialising config: %w", err)
	}

	if err := w.initPages(); err != nil {
		return fmt.Errorf("Error initialising pages: %w", err)
	}

	if err := w.initLayouts(); err != nil {
		return fmt.Errorf("Error initialising layouts: %w", err)
	}

	if err := w.initIncludes(); err != nil {
		return fmt.Errorf("Error initialising includes: %w", err)
	}

	if err := w.initPosts(); err != nil {
		return fmt.Errorf("Error initialising posts: %w", err)
	}

	return nil
}

func (w *Website) initConfig() error {
	w.Config = &Config{}

	f := "_config.yml"
	p := fmt.Sprintf("%s/%s", w.SourcePath, f)
	if err := w.Config.SetFromFile(p); err != nil {
		return fmt.Errorf("Error setting config from %s: %w", p, err)
	}

	if err := w.Config.Validate(); err != nil {
		return fmt.Errorf("Config is invalid: %w", err)
	}

	return nil
}

func (w *Website) initPages() error {
	w.PageNames = []string{}
	w.Pages = map[string]*Page{}

	entries, err := os.ReadDir(w.SourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Source directory does not exist")
		}
		return fmt.Errorf("Error reading source directory: %w", err)
	}

	for _, e := range entries {
		entryPath := filepath.Join(w.SourcePath, e.Name())

		fileInfo, err := os.Stat(entryPath)
		if err != nil {
			continue
		}
		re := regexp.MustCompile(`^[a-zA-Z0-9\_\-]+\.(markdown|html)$`)
		if !fileInfo.Mode().IsRegular() || !re.MatchString(e.Name()) {
			continue
		}

		page := &Page{}
		if err := page.SetFromFile(entryPath); err != nil {
			return fmt.Errorf("Error getting page from %s: %w", entryPath, err)
		}

		if err := page.Validate(); err != nil {
			return fmt.Errorf("Page %s is invalid: %w", page.Name, err)
		}

		w.PageNames = append(w.PageNames, page.Name)
		w.Pages[page.Name] = page
	}

	if w.Pages["index"] == nil {
		return errors.New("Cannot find index.markdown")
	}

	return nil
}

func (w *Website) initLayouts() error {
	w.LayoutNames = []string{}
	w.Layouts = map[string]*Layout{}

	names, err := w.getFilenamesWithExtensionsFromDir("_layouts")
	if err != nil {
		return fmt.Errorf("Error getting layouts: %w", err)
	}

	for _, n := range names {
		w.LayoutNames = append(w.LayoutNames, n)
		w.Layouts[n] = &Layout{
			Name: n,
		}

		p := filepath.Join(w.SourcePath, "_layouts", n+".html")
		if err := w.Layouts[n].SetFromFile(p); err != nil {
			return fmt.Errorf("Error setting layout from %s: %w", p, err)
		}
	}

	return nil
}

func (w *Website) initIncludes() error {
	w.IncludeNames = []string{}
	w.Includes = map[string]*Include{}

	names, err := w.getFilenamesWithExtensionsFromDir("_includes")
	if err != nil {
		return fmt.Errorf("Error getting layouts: %w", err)
	}

	for _, n := range names {
		w.IncludeNames = append(w.IncludeNames, n)
		w.Includes[n] = &Include{
			Name: n,
		}

		p := filepath.Join(w.SourcePath, "_includes", n+".html")
		if err := w.Includes[n].SetFromFile(p); err != nil {
			return fmt.Errorf("Error setting include from %s: %w", p, err)
		}
	}

	return nil
}

func (w *Website) initPosts() error {
	w.PostsNames = []string{}
	w.Posts = map[string]*Page{}

	names, err := w.getFilenamesWithExtensionsFromDir("_posts")
	if err != nil {
		return fmt.Errorf("Error getting posts: %w", err)
	}

	for _, n := range names {
		w.PostsNames = append(w.PostsNames, n)
		w.Posts[n] = &Page{
			Name: n,
		}

		p := filepath.Join(w.SourcePath, "_posts", n+".markdown")
		if err := w.Posts[n].SetFromFile(p); err != nil {
			return fmt.Errorf("Error setting post from %s: %w", p, err)
		}
	}

	return nil
}

func (w *Website) getFilenamesWithExtensionsFromDir(d string) ([]string, error) {
	p := fmt.Sprintf("%s/%s", w.SourcePath, d)

	entries, err := os.ReadDir(p)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, fmt.Errorf("Directory %s does not exist", p)
		}
		return []string{}, fmt.Errorf("Error reading directory %s: %w", p, err)
	}

	names := []string{}
	foundNames := map[string]bool{}
	re := regexp.MustCompile(`^[a-zA-Z0-9\_\-]+\.(html|markdown)$`)
	for _, e := range entries {
		entryPath := filepath.Join(p, e.Name())
		fileInfo, err := os.Stat(entryPath)
		if err != nil {
			return []string{}, fmt.Errorf("Error getting file info for %s: %s", entryPath, err)
		}
		if !fileInfo.Mode().IsRegular() || !re.MatchString(e.Name()) {
			continue
		}

		nameWithoutExt := strings.Replace(e.Name(), ".html", "", -1)
		nameWithoutExt = strings.Replace(nameWithoutExt, ".markdown", "", -1)

		if foundNames[nameWithoutExt] {
			return []string{}, fmt.Errorf("There are two files %s but with different extensions", nameWithoutExt)
		}

		foundNames[nameWithoutExt] = true
		names = append(names, nameWithoutExt)
	}

	return names, nil
}
