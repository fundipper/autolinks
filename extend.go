package autolinks

import (
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

type Extender struct {
	Data map[string]string
}

// New return initialized image render with source url replacing support.
func NewExtender(data map[string]string) goldmark.Extender {
	return &Extender{
		Data: data,
	}
}

func (e *Extender) Extend(m goldmark.Markdown) {
	ps := []util.PrioritizedValue{}
	for k, v := range e.Data {
		ps = append(ps, util.Prioritized(
			NewTransformer(regexp.MustCompile(k), []byte(v)), 500),
		)
	}

	m.Parser().AddOptions(
		parser.WithASTTransformers(ps...),
	)
	// m.Renderer().AddOptions(
	// 	renderer.WithNodeRenderers(
	// 		util.Prioritized(NewRenderer(), 500),
	// 	),
	// )
}
