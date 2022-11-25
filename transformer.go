package autolinks

import (
	"log"
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Transformer struct {
	Pattern *regexp.Regexp
	URL     []byte
}

func NewTransformer(pattern *regexp.Regexp, url []byte) *Transformer {
	return &Transformer{
		Pattern: pattern,
		URL:     url,
	}
}

// Transform implements goldmark.parser.ASTTransformer
func (t *Transformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	// Walk the AST in depth-first fashion and apply transformations
	err := ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		// Each node will be visited twice, once when it is first encountered (entering), and again
		// after all the node's children have been visited (if any). Skip the latter.
		if !entering {
			return ast.WalkContinue, nil
		}
		// Skip the children of existing links to prevent double-transformation.
		if node.Kind() == ast.KindLink || node.Kind() == ast.KindAutoLink {
			return ast.WalkSkipChildren, nil
		}
		// Linkify any Text nodes encountered
		if node.Kind() == ast.KindText {
			textNode := node.(*ast.Text)
			t.LinkifyText(textNode, reader.Source())
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		log.Fatal("Error encountered while transforming AST:", err)
	}
}

// LinkifyText finds all LinkPattern matches in the given Text node and replaces them with Link
// nodes that point to ReplUrl.
func (t *Transformer) LinkifyText(node *ast.Text, source []byte) {
	parent := node.Parent()
	tSegment := node.Segment

	match := t.Pattern.FindIndex(tSegment.Value(source))
	if match == nil {
		return
	}
	// Create a text.Segment for the link text.
	lSegment := text.NewSegment(tSegment.Start+match[0], tSegment.Start+match[1])

	// Insert node for any text before the link
	if lSegment.Start != tSegment.Start {
		bText := ast.NewTextSegment(tSegment.WithStop(lSegment.Start))
		parent.InsertBefore(parent, node, bText)
	}

	// Insert Link node
	link := ast.NewLink()
	link.AppendChild(link, ast.NewTextSegment(lSegment))
	link.Destination = t.Pattern.ReplaceAll(lSegment.Value(source), t.URL)
	parent.InsertBefore(parent, node, link)

	// Update original node to represent the text after the link (may be empty)
	node.Segment = tSegment.WithStart(lSegment.Stop)

	// Linkify remaining text if not empty
	if node.Segment.Len() > 0 {
		t.LinkifyText(node, source)
	}
}
