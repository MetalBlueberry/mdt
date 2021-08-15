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
			if entering {
				return ast.WalkContinue, nil
			}
			if string(tnode.Language(source)) == "mermaid" {

				previous := tnode.PreviousSibling()
				previousBlock, hasPreviousblock := previous.(*ast.HTMLBlock)

				next := tnode.NextSibling()
				nextBlock, hasNextBlock := next.(*ast.HTMLBlock)

				if hasPreviousblock && hasNextBlock {
					start := previousBlock.Lines().At(0).Start
					stop := nextBlock.Lines().At(next.Lines().Len() - 1).Stop
					_, err := parseHTML(source[start:stop])
					if err == nil {
						return ast.WalkContinue, nil
					}
				}

				segment := slice(tnode.Lines())
				code := source[segment.Start:segment.Stop]
				segment.Start -= len("```mermaid") + 1
				segment.Stop += len("```") + 1
				fences = append(fences, &Fence{
					Segment: segment,
					Code:    code,
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
