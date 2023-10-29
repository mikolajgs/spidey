package main

import (
	"fmt"
	"os"
	"strings"
	"regexp"
)

type Node struct {
	Type string
	Content string
	Children []*Node
	Parent *Node
	Values map[string]map[string]string
}

func (n *Node) ProcessRawTags(tagPrefix string, tagSuffix string) {
	if n.Type == "raw" {
		n.Type = "text"
		c := ""
		if len(n.Children) > 0 {
			for _, child := range n.Children {
				c += child.GetRaw(tagPrefix, tagSuffix)
				child = nil
			}
			n.Children = []*Node{}
		}
		n.Content = c
	} else {
		if len(n.Children) > 0 {
			for _, child := range n.Children {
				child.ProcessRawTags(tagPrefix, tagSuffix)
			}
		}
	}
}

func (n *Node) GetRaw(tagPrefix string, tagSuffix string) string {
	s := ""
	if n.Type == "text" {
		s += n.Content
	} else if n.Type == "if" || n.Type == "for" || n.Type == "raw" {
		s += fmt.Sprintf("%s%s%s", tagPrefix, n.Content, tagSuffix)
	}
	if len(n.Children) > 0 {
		for _, child := range n.Children {
			s += child.GetRaw(tagPrefix, tagSuffix)
		}
	}
	if n.Type == "if" || n.Type == "for" || n.Type == "raw" {
		s += fmt.Sprintf("%send%s%s", tagPrefix, n.Type, tagSuffix)
	}
	return s
}

func (n *Node) ProcessForTags(tagPrefix string, tagSuffix string, w *Website, g *Generator) {
	// TODO: Hardcoded!
	if n.Type == "for" && strings.Trim(n.Content, " ") == "for post in site.posts" {
		forContent := ""
		for _, ch := range n.Children {
			forContent += ch.GetRaw(tagPrefix, tagSuffix)
		}

		newChildren := []*Node{}
		for _, post := range w.Posts {
			childNode := &Node {
				Type: "group",
				Values: map[string]map[string]string{},
			}
			childNode.Values["post"] = g.getObjVariablesFromYamlTag(post)
			childNode.SetFromString(forContent, []rune(tagPrefix)[0], []rune(tagSuffix)[1], []rune(tagPrefix)[1])
			newChildren = append(newChildren, childNode)
		}
		n.Children = newChildren
		n.Type = "group"
	}
	if len(n.Children) > 0 {
		for _, child := range n.Children {
			child.ProcessForTags(tagPrefix, tagSuffix, w, g)
		}
	}
}

func (n *Node) ProcessIfTags(siteVars map[string]string, pageVars map[string]string) {
	if n.Type == "if" {
		s := strings.Trim(n.Content, " ")
		sArr := strings.Split(s, " ")
		result := false
		if len(sArr) > 2 || sArr[0] != "if" {
			n.Type = "text"
			n.Content = "INVALID IF"
			n.Children = []*Node{}
		}

		re := regexp.MustCompile(`(site|page|post)\.([a-zA-Z0-9\-\_]+)`)
		found := re.FindStringSubmatch(sArr[1])
		if len(found) != 3 {
			n.Type = "Text"
			n.Content = "INVALID IF"
			n.Children = []*Node{}
		}

		if found[1] == "site" {
			if siteVars[found[2]] != "" {
				result = true
			}
		} else if found[1] == "page" {
			if pageVars[found[2]] != "" {
				result = true
			}
		} else if found[1] == "post" {
			if n.GetNodeAttachedValue("post", found[2]) != "" {
				result = true
			}
		}

		if !result {
			n.Type = "text"
			n.Content = ""
			n.Children = []*Node{}
		} else {
			n.Type = "group"
			n.Content = ""
		}
	}
	if len(n.Children) > 0 {
		for _, child := range n.Children {
			child.ProcessIfTags(siteVars, pageVars)
		}
	}
}

func (n *Node) ProcessPostVars() {
	re := regexp.MustCompile(`\{\{[ ]*post\.([a-zA-Z0-9\-\_]+)[ ]*\}\}`)
	if n.Type == "text" {
		for _, found := range re.FindAllStringSubmatch(n.Content, -1) {
			varName := strings.Trim(string(found[1]), " ")
			replaceWith := n.GetNodeAttachedValue("post", varName)
			n.Content = strings.ReplaceAll(n.Content, found[0], replaceWith)
		}
	}
	for _, ch := range n.Children {
		ch.ProcessPostVars()
	}
}

func (n *Node) GetNodeAttachedValue(objName string, varName string) string {
	if n.Values[objName] != nil {
		return n.Values[objName][varName]
	} else if n.Type != "root" {
		return n.Parent.GetNodeAttachedValue(objName, varName)
	}
	return ""
}

func (n *Node) SetFromString(h string, openRune rune, closeRune rune, tagRune rune) {
	var prevCh rune
	var prevPrevCh rune
	tagStarted := false
	tagContents := ""
	text := []rune("")

	lastNode := n

	runes := []rune(h)

	for i, ch := range runes {
		if ch == tagRune && prevCh == openRune && !tagStarted {
			tagStarted = true
			prevPrevCh = prevCh
			prevCh = ch
			continue
		} 

		if ch == closeRune && prevCh == tagRune && tagStarted {
			tagStarted = false

			tagName := n.getTagName(tagContents, openRune, closeRune, tagRune)
			if tagName == "if" || tagName == "for" || tagName == "raw" {
				node1 := &Node{
					Type: "text",
					Content: string(text),
					Parent: lastNode,
				}
				lastNode.Children = append(lastNode.Children, node1)
				node2 := &Node{
					Type: tagName,
					Content: tagContents,
					Children: []*Node{},
					Parent: lastNode,
				}
				lastNode.Children = append(lastNode.Children, node2)
				lastNode = node2
			}
			if tagName == "endif" || tagName == "endfor" || tagName == "endraw" {
				node := &Node{
					Type: "text",
					Content: string(text),
					Parent: lastNode,
				}
				lastNode.Children = append(lastNode.Children, node)
				lastNode = lastNode.Parent
			}

			tagContents = ""
			prevPrevCh = prevCh
			prevCh = ch
			text = []rune("")
			continue
		}

		if tagStarted && prevCh != tagRune {
			tagContents += string(prevCh)
			prevPrevCh = prevCh
			prevCh = ch
			continue
		}

		if !tagStarted && !(prevCh == closeRune && prevPrevCh == tagRune) {
			if i > 0 && i != len(runes)-1 {
				text = append(text, prevCh)
			}
			if i == len(runes)-1 {
				text = append(text, prevCh)
				text = append(text, ch)

				node := &Node{
					Type: "text",
					Content: string(text),
					Parent: lastNode,
				}
				lastNode.Children = append(lastNode.Children, node)
			}
		}

		prevPrevCh = prevCh
		prevCh = ch
	}
}

func (n *Node) getTagName(s string, openRune rune, closeRune rune, tagRune rune) string {
	s = strings.Replace(s, string([]rune{openRune, tagRune}), "", -1)
	s = strings.Replace(s, string([]rune{tagRune, closeRune}), "", -1)
	s = strings.Trim(s, " ")
	sArr := strings.Split(s, " ")
	return sArr[0]
}

func (n *Node) Debug(depth int) {
	fmt.Fprintf(os.Stdout, "%sNode: type=%s", strings.Repeat(" ", depth*2), n.Type)
	if n.Parent != nil {
		fmt.Fprintf(os.Stdout, " parent_type=%s", n.Parent.Type)
	}
	content := strings.ReplaceAll(n.Content, "\n", "\\n")
	if len(n.Content) > 180 {
		content = "TOO LONG"
	}
	fmt.Fprintf(os.Stdout, " content=%s", content)
	fmt.Fprintf(os.Stdout, "\n")
	if len(n.Children) > 0 {
		for _, c := range n.Children {
			c.Debug(depth+1)
		}
	}
}
