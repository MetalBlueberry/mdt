package mdt

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Fence struct {
	Segment text.Segment
	Code    []byte
}

func ParseFences(source []byte, root ast.Node) ([]*Fence, error) {
	fences := []*Fence{}
	err := ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch tnode := n.(type) {
		case *ast.FencedCodeBlock:
			if !entering {
				return ast.WalkContinue, nil
			}
			if string(tnode.Language(source)) == "mermaid" {
				segment := slice(tnode.Lines())
				fences = append(fences, &Fence{
					Segment: segment,
					Code:    source[segment.Start:segment.Stop],
				})
			}
		default:
			// fmt.Printf("%T : %s : %t\n", n, string(n.Text(source)), n.IsRaw())
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}
	return fences, nil
}

func (f *Fence) Slice() text.Segment {
	return f.Segment
}

func (f *Fence) Content() []byte {
	return f.Code
}

func slice(lines *text.Segments) text.Segment {
	lastLine := lines.Len() - 1
	first := lines.At(0).Start
	last := lines.At(lastLine).Stop
	return text.NewSegment(first, last)
}
