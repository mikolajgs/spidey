package main

import (
	"testing"
)

func TestGetRaw(t *testing.T) {
	s := testNode1.GetRaw("{%", "%}")
	if s != "NodeStart!Fruits?{%if page.fruit%}Apple.{%raw%}{%if page.sweet%}Sweet ones.{%if page.yellow%}Banana.Yellow.{%endif%}{%endif%}{%endraw%}{%endif%}" {
		t.Fatalf("GetRaw returns invalid string")
	}
}

func TestProcessRawTags(t *testing.T) {
	testNode1.ProcessRawTags("{%", "%}")

	if len(testNode1.Children) != 3 {
		t.Fatalf("ProcessRawTags returns invalid number of children")
	}

	if testNode1.Children[2].Type != "group" ||
		len(testNode1.Children[2].Children) != 1 ||
		testNode1.Children[2].Children[0].Type != "if" ||
		len(testNode1.Children[2].Children[0].Children) != 2 ||
		testNode1.Children[2].Children[0].Children[0].Type != "text" ||
		testNode1.Children[2].Children[0].Children[1].Type != "text" ||
		testNode1.Children[2].Children[0].Children[1].Content != "{%if page.sweet%}Sweet ones.{%if page.yellow%}Banana.Yellow.{%endif%}{%endif%}" {

		t.Fatalf("ProcessRawTags failed to process raw tags")
	}
}

func TestProcessForTags(t *testing.T) {
	testNode2.ProcessForTags("{%", "%}", testWebsite2, testGenerator2)
	s := testNode2.GetRaw("{%", "%}")
	if s != "Post: {{ post.title }}{%if post.description%}Description: {{ post.description }}{%endif%}--Post: {{ post.title }}{%if post.description%}Description: {{ post.description }}{%endif%}--" {
		t.Fatalf("ProcessForTags failed to process 'for' tag")
	}
}

func TestProcessIfTags(t *testing.T) {
	siteVars := map[string]string{
		"title": "SiteTitle",
	}
	pageVars := map[string]string{
		"title": "PageTitle",
	}

	testNode3.ProcessIfTags(siteVars, pageVars)
	s := testNode3.GetRaw("{%", "%}")
	if s != "TestIf!PageTitle!TestIfAgain!SiteTitle!" {
		t.Fatalf("ProcessIfTags failed to process pages and sites in 'if' tags")
	}

	testNode2.ProcessIfTags(siteVars, pageVars)
	s = testNode2.GetRaw("{%", "%}")
	if s != "Post: {{ post.title }}Description: {{ post.description }}--Post: {{ post.title }}--" {
		t.Fatalf("ProcessForTags failed to process posts in 'if' tags")
	}
}

func TestProcessPostVars(t *testing.T) {
	testNode2.ProcessPostVars()
	s := testNode2.GetRaw("{%", "%}")
	if s != "Post: Title1Description: Description1--Post: Title2--" {
		t.Fatalf("ProcessPostVars failed to process post values")
	}
}

func TestGetNodeAttachedValue(t *testing.T) {
	testNode2b.ProcessForTags("{%", "%}", testWebsite2, testGenerator2)
	t1 := testNode2b.Children[0].Children[0].Children[1].Children[0].GetNodeAttachedValue("post", "title")
	t2 := testNode2b.Children[0].Children[1].Children[1].Children[0].GetNodeAttachedValue("post", "title")
	d1 := testNode2b.Children[0].Children[0].Children[1].Children[0].GetNodeAttachedValue("post", "description")
	d2 := testNode2b.Children[0].Children[1].Children[1].Children[0].GetNodeAttachedValue("post", "description")
	if t1 != "Title1" ||
		t2 != "Title2" ||
		d1 != "Description1" ||
		d2 != "" {

		t.Fatalf("GetNodeAttachedValue failed to get attached value")
	}
}

func TestSetFromString(t *testing.T) {
	n := &Node{}
	n.SetFromString(""+
		"TextBlock!"+
		"{%if page.var1%}"+
		"TextVar1!"+
		"{%if page.var2%}"+
		"TextVar2!"+
		"{%endif%}"+
		"StillTextVar1!"+
		"{%endif%}"+
		"Text!"+
		"{%raw%}"+
		"Raw!"+
		"{%if page.var3%}"+
		"PageVar3!"+
		"{%endif%}"+
		"{%endraw%}"+
		"Text2!"+
		"{%for post in site.posts%}"+
		"{%if post.var1%}"+
		"PostVar1!"+
		"{%endif%}"+
		"{%endfor%}"+
		"Text3!", '{', '}', '%')

	if len(n.Children) != 7 ||
		len(n.Children[1].Children) != 3 ||
		len(n.Children[3].Children) != 3 ||
		len(n.Children[5].Children) != 3 {

		t.Fatalf("SetFromString failed to create valid number of children nodes")
	}
	for i, v := range []string{"text", "if", "text", "raw", "text", "for", "text"} {
		if n.Children[i].Type != v {
			t.Fatalf("SetFromString failed to create children of valid types")
		}
	}

	if n.Children[1].Children[1].Type != "if" ||
		n.Children[1].Children[1].Content != "if page.var2" ||
		len(n.Children[1].Children[1].Children) != 1 ||
		n.Children[1].Children[1].Children[0].Type != "text" ||
		n.Children[1].Children[1].Children[0].Content != "TextVar2!" {

		t.Fatalf("SetFromString failed to parse 'if' tag")
	}

	if n.Children[3].Children[1].Type != "if" ||
		n.Children[3].Children[1].Content != "if page.var3" ||
		len(n.Children[3].Children[1].Children) != 1 ||
		n.Children[3].Children[1].Children[0].Type != "text" ||
		n.Children[3].Children[1].Children[0].Content != "PageVar3!" {

		t.Fatalf("SetFromString failed to parse 'raw' tag")
	}

	if len(n.Children[5].Children) != 3 ||
		n.Children[5].Children[0].Content != "" ||
		n.Children[5].Children[2].Content != "" ||
		n.Children[5].Children[1].Type != "if" ||
		n.Children[5].Children[1].Content != "if post.var1" ||
		len(n.Children[5].Children[1].Children) != 1 ||
		n.Children[5].Children[1].Children[0].Type != "text" ||
		n.Children[5].Children[1].Children[0].Content != "PostVar1!" {

		t.Fatalf("SetFromString failed to parse 'for' tag")
	}
	if n.Children[0].Content != "TextBlock!" ||
		n.Children[2].Content != "Text!" ||
		n.Children[4].Content != "Text2!" ||
		n.Children[6].Content != "Text3!" {
		t.Fatalf("SetFromString failed to parse text blocks")
	}
}
