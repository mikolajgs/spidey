package main

import (
	"errors"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type Generator struct {
	DestinationPath string

	cachedSiteVariables map[string]string
}

func (g *Generator) Generate(w *Website) error {
	err := g.checkIfDestinationPathEmpty()
	if err != nil {
		return err
	}

	g.getSiteVariables(w.Config)

	if err := g.generatePosts(w); err != nil {
		return err
	}

	if err := g.generatePages(w); err != nil {
		return err
	}

	return nil
}

func (g *Generator) checkIfDestinationPathEmpty() error {
	fileInfo, err := os.Stat(g.DestinationPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Destination path %s does not exist", g.DestinationPath)
		}
		if !fileInfo.IsDir() {
			return fmt.Errorf("Destination path %s is not a directory", g.DestinationPath)
		}
		return fmt.Errorf("Error getting file info for %s: %w", g.DestinationPath, err)
	}

	entries, err := os.ReadDir(g.DestinationPath)
	if err != nil {
		return fmt.Errorf("Error reading destination path of %s: %w", g.DestinationPath, err)
	}

	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), ".") {
			return fmt.Errorf("Destination path %s is not empty.  It can only contain dot files", g.DestinationPath)
		}
	}

	return nil
}

func (g *Generator) getSiteVariables(c *Config) {
	g.cachedSiteVariables = g.getObjVariablesFromYamlTag(c)
}

func (g *Generator) getObjVariablesFromYamlTag(obj interface{}) map[string]string {
	out := map[string]string{}

	v := reflect.ValueOf(obj)
	i := reflect.Indirect(v)
	s := i.Type()
	for j := 0; j < s.NumField(); j++ {
		field := s.Field(j)
		fieldKind := field.Type.Kind()
		if fieldKind != reflect.String {
			continue
		}
		tagVal := field.Tag.Get("yaml")
		if tagVal == "" {
			continue
		}
		tagValArr := strings.Split(tagVal, ",")
		out[tagValArr[0]] = v.Elem().FieldByName(field.Name).String()
	}
	return out
}

func (g *Generator) generatePages(w *Website) error {
	for name, page := range w.Pages {
		pageHtml, err := g.getPageHtml(page, w)
		if err != nil {
			return fmt.Errorf("Error generating page %s HTML: %w", name, err)
		}

		pagePath := filepath.Join(g.DestinationPath, "index.html")
		if name != "index" && name != "404" {
			pageDir := filepath.Join(g.DestinationPath, name)
			err := os.Mkdir(pageDir, 0750)
			if err != nil {
				return fmt.Errorf("Error creating page %s directory: %w", name, err)
			}
			pagePath = filepath.Join(pageDir, "index.html")
		}
		if name == "404" {
			pagePath = filepath.Join(g.DestinationPath, "404.html")
		}

		err = os.WriteFile(pagePath, []byte(pageHtml), 0750)
		if err != nil {
			return fmt.Errorf("Error writing page html to %s: %w", pagePath, err)
		}
	}
	return nil
}

func (g *Generator) generatePosts(w *Website) error {
	for name, post := range w.Posts {
		postHtml, err := g.getPageHtml(post, w)
		if err != nil {
			return fmt.Errorf("Error generating post %s HTML: %w", name, err)
		}

		re := regexp.MustCompile(`^([0-9]{4})-([01][0-9])-([01][0-9])-([a-zA-Z0-9\_\-]+)$`)
		if !re.MatchString(name) {
			return fmt.Errorf("Post %s filename does not match regexp", name)
		}
		nameArr := re.FindStringSubmatch(name)

		destPath := []string{}
		re = regexp.MustCompile(`^[a-zA-Z0-9\_\- ]+$`)
		if post.Categories != "" && re.MatchString(post.Categories) {
			categoriesList := strings.Split(post.Categories, " ")
			for _, c := range categoriesList {
				if c != "" {
					destPath = append(destPath, c)
				}
			}
		} else {
			destPath = append(destPath, "posts")
		}

		postDir := filepath.Join(destPath...)
		postDir = filepath.Join(postDir, nameArr[1], nameArr[2], nameArr[3])
		err = os.MkdirAll(filepath.Join(g.DestinationPath, postDir), 0750)
		if err != nil {
			return fmt.Errorf("Error creating post %s directory %s: %w", name, postDir, err)
		}

		postPath := filepath.Join(postDir, "index.html")
		w.Posts[name].Url = "/" + postPath

		err = os.WriteFile(filepath.Join(g.DestinationPath, postPath), []byte(postHtml), 0750)
		if err != nil {
			return fmt.Errorf("Error writing post html to %s: %w", postPath, err)
		}
	}
	return nil
}

func (g *Generator) getPageHtml(p *Page, w *Website) (string, error) {
	if w.Layouts[p.Layout] == nil {
		return "", fmt.Errorf("Layout %s does not exist", p.Layout)
	}

	pageHtml := w.Layouts[p.Layout].Body

	contentHtml := ""
	if p.ContentType == "html" {
		contentHtml = p.Body
	} else if p.ContentType == "markdown" {
		contentHtml = g.mdToHtml(p.Body)
	}

	re := regexp.MustCompile(`\{\{[ ]*content[ ]*\}\}`)
	for _, cont := range re.FindAllStringSubmatch(pageHtml, -1) {
		pageHtml = strings.ReplaceAll(pageHtml, cont[0], contentHtml)
	}

	// TODO: Fix it to actually replace until necessary, with either limit or infinite
	// cycle detection
	var err error
	for i := 0; i < 10; i++ {
		pageHtml, err = g.replaceIncludes(pageHtml, w)
		if err != nil {
			return "", fmt.Errorf("Error replacing includes in page %s: %w", p.Name, err)
		}
	}

	pageHtml, err = g.processTags(pageHtml, w, p)
	pageHtml = g.addBaseUrl(pageHtml, w)

	return pageHtml, nil
}

func (g *Generator) processTags(s string, w *Website, p *Page) (string, error) {
	s, err := g.replaceOnTree(s, w, p)
	if err != nil {
		return "", fmt.Errorf("Error replacing ifs and fors: %w", err)
	}
	s, err = g.replaceVariables(s, w, p)
	if err != nil {
		return "", fmt.Errorf("Error replacing vars: %w", err)
	}
	return s, nil
}

func (g *Generator) mdToHtml(md string) string {
	replaced := map[string]string{}
	re := regexp.MustCompile(`\{%[ ]*include[ ]*([a-zA-Z0-9\-\_]+)\.(html|markdown)[ ]*%\}`)
	for _, incl := range re.FindAllStringSubmatch(md, -1) {
		replacement := fmt.Sprintf("<!--- TMPTAG:%s -->", incl[0])
		replaced[incl[0]] = replacement
		md = strings.ReplaceAll(md, incl[0], replacement)
	}

	re = regexp.MustCompile(`\{%[ ]*(endraw|raw)[ ]*%\}`)
	for _, tags := range re.FindAllStringSubmatch(md, -1) {
		replacement := fmt.Sprintf("<!--- TMPTAG:%s -->", tags[0])
		replaced[tags[0]] = replacement
		md = strings.ReplaceAll(md, tags[0], replacement)
	}

	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	h := string(markdown.Render(doc, renderer))
	for k, v := range replaced {
		h = strings.ReplaceAll(h, v, k)
	}

	return h
}

func (g *Generator) replaceIncludes(h string, w *Website) (string, error) {
	re := regexp.MustCompile(`\{%[ ]*include[ ]*([a-zA-Z0-9\-\_]+)\.(html|markdown)[ ]*%\}`)
	for _, incl := range re.FindAllStringSubmatch(h, -1) {
		inclName := incl[1]
		if w.Includes[inclName] == nil {
			return "", fmt.Errorf("Include %s does not exist", inclName)
		}
		h = strings.ReplaceAll(h, incl[0], w.Includes[inclName].Body)
	}
	return h, nil
}

func (g *Generator) replaceVariables(h string, w *Website, p *Page) (string, error) {
	pageVars := g.getObjVariablesFromYamlTag(p)
	re := regexp.MustCompile(`\{\{[ ]*(site|page)\.([a-zA-Z0-9\-\_]+)[ ]*\}\}`)
	for _, found := range re.FindAllStringSubmatch(h, -1) {
		varType := found[1]
		varName := found[2]
		if varType == "page" {
			h = strings.ReplaceAll(h, found[0], pageVars[varName])
		} else if varType == "site" {
			h = strings.ReplaceAll(h, found[0], g.cachedSiteVariables[varName])
		} else {
			return "", errors.New(fmt.Sprintf("Invalid variable type %s", varType))
		}
	}
	return h, nil
}

func (g *Generator) addBaseUrl(h string, w *Website) string {
	re := regexp.MustCompile(`href="/`)
	url := w.Config.Url
	if w.Config.Baseurl != "" {
		url = fmt.Sprintf("%s/%s", url, w.Config.Baseurl)
	}
	for _, href := range re.FindAllStringSubmatch(h, -1) {
		h = strings.ReplaceAll(h, href[0], fmt.Sprintf("href=\"%s/", url))
	}
	return h
}

func (g *Generator) replaceOnTree(h string, w *Website, p *Page) (string, error) {
	tree := &Node{
		Type: "root",
	}
	tree.SetFromString(h, '{', '}', '%')
	tree.ProcessRawTags("{%", "%}")
	tree.ProcessForTags("{%", "%}", w, g)

	pageVars := g.getObjVariablesFromYamlTag(p)
	tree.ProcessIfTags(g.cachedSiteVariables, pageVars)
	tree.ProcessPostVars()

	h = tree.GetRaw("", "")

	return h, nil
}
