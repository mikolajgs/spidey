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

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
