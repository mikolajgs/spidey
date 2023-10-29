package main

import (
	"os"
	"testing"
)

var testNode1 = &Node{
	Type: "root",
	Children: []*Node{
		&Node{
			Type: "text",
			Content: "NodeStart!",
		},
		&Node{
			Type: "text",
			Content: "Fruits?",
		},
		&Node{
			Type: "group",
			Children: []*Node{
				&Node{
					Type: "if",
					Content: "if page.fruit",
					Children: []*Node{
						&Node{
							Type: "text",
							Content: "Apple.",
						},
						&Node{
							Type: "raw",
							Content: "raw",
							Children: []*Node{
								&Node{
									Type: "if",
									Content: "if page.sweet",
									Children: []*Node{
										&Node{
											Type: "text",
											Content: "Sweet ones.",
										},
										&Node{
											Type: "if",
											Content: "if page.yellow",
											Children: []*Node{
												&Node{
													Type: "text",
													Content: "Banana.",
												},
												&Node{
													Type: "text",
													Content: "Yellow.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}


var testNode2 = &Node{
	Type: "root",
	Children: []*Node{
		&Node{
			Type: "for",
			Content: "for post in site.posts",
			Children: []*Node{
				&Node{
					Type: "text",
					Content: "Post: {{ post.title }}",
				},
				&Node{
					Type: "if",
					Content: "if post.description",
					Children: []*Node{
						&Node{
							Type: "text",
							Content: "Description: {{ post.description }}",
						},
					},
				},
				&Node{
					Type: "text",
					Content: "--",
				},
			},
		},
	},
}

var testNode2b = &Node{
	Type: "root",
	Children: []*Node{
		&Node{
			Type: "for",
			Content: "for post in site.posts",
			Children: []*Node{
				&Node{
					Type: "text",
					Content: "Post: {{ post.title }}",
				},
				&Node{
					Type: "if",
					Content: "if post.description",
					Children: []*Node{
						&Node{
							Type: "text",
							Content: "Description: {{ post.description }}",
						},
					},
				},
				&Node{
					Type: "text",
					Content: "--",
				},
			},
		},
	},
}

var testWebsite2 = &Website{
	Posts: map[string]*Page{
		"one": &Page{
			Title: "Title1",
			Description: "Description1",
		},
		"two": &Page{
			Title: "Title2",
			Description: "",
		},
	},
}

var testGenerator2 = &Generator{}


var testNode3 = &Node{
	Type: "root",
	Children: []*Node{
		&Node{
			Type: "text",
			Content: "TestIf!",
			Children: []*Node{
				&Node{
					Type: "if",
					Content: "if page.title",
					Children: []*Node{
						&Node{
							Type: "text",
							Content: "PageTitle!",
						},
					},
				},
				&Node{
					Type: "if",
					Content: "if page.description",
					Children: []*Node{
						&Node{
							Type: "text",
							Content: "PageDescription!",
						},
					},
				},
			},
		},
		&Node{
			Type: "text",
			Content: "TestIfAgain!",
			Children: []*Node{
				&Node{
					Type: "if",
					Content: "if site.title",
					Children: []*Node{
						&Node{
							Type: "text",
							Content: "SiteTitle!",
						},
					},
				},
				&Node{
					Type: "if",
					Content: "if site.description",
					Children: []*Node{
						&Node{
							Type: "text",
							Content: "SiteDescription!",
						},
					},
				},
			},
		},
	},
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
