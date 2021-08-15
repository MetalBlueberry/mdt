package mdt

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type FenceWrap struct {
	Segment text.Segment
	Block   Block
}

func ParseWrappedFences(source []byte, root ast.Node) ([]*FenceWrap, error) {
	fencesWrap := []*FenceWrap{}
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

				if !hasPreviousblock || !hasNextBlock {
					return ast.WalkContinue, nil
				}

				start := previousBlock.Lines().At(0).Start
				stop := nextBlock.Lines().At(next.Lines().Len() - 1).Stop

				block, err := parseHTML(source[start:stop])
				if err != nil {
					return ast.WalkContinue, err
				}
				segment := slice(tnode.Lines())
				block.Code.Mermaid = source[segment.Start:segment.Stop]

				fencesWrap = append(fencesWrap, &FenceWrap{
					Segment: text.NewSegment(start, stop),
					Block:   block,
				})
			}
		// case *ast.HTMLBlock:
		// 	if !entering {
		// 		return ast.WalkContinue, nil
		// 	}
		// 	lastLine := tnode.Lines().Len() - 1
		// 	first := tnode.Lines().At(0).Start
		// 	last := tnode.Lines().At(lastLine).Stop

		// 	block, err := parseHTML(source[first:last])
		// 	if err != nil {
		// 		return ast.WalkStop, err
		// 	}
		// 	fencesWrap = append(fencesWrap, &FenceWrap{
		// 		Segment: text.NewSegment(first, last),
		// 		Block:   block,
		// 	})

		default:
			// fmt.Printf("%T : %s : %t\n", n, string(n.Text(source)), n.IsRaw())
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}
	return fencesWrap, nil
}

func NewBlock(code []byte, src string) Block {
	return Block{
		XMLName: xml.Name{Local: "details"},
		Class:   "mermaid",
		Img: Img{
			Src: src,
		},
		Code: Code{
			Code:    fmt.Sprintf("\n\n```mermaid\n%s```\n", string(code)),
			Mermaid: code,
		},
	}
}

type Block struct {
	XMLName xml.Name
	Class   string `xml:"class,attr"`
	Img     Img    `xml:"summary>img"`
	Code    Code   `xml:"p"`
}

type Img struct {
	Src string `xml:"src,attr"`
}
type Code struct {
	Code    string `xml:",innerxml"`
	Mermaid []byte `xml:"-"`
}

func parseHTML(src []byte) (Block, error) {
	block := Block{}

	decoder := xml.NewDecoder(bytes.NewReader(src))

	err := decoder.Decode(&block)
	if err != nil {
		return Block{}, err
	}
	if block.Class != "mermaid" {
		return block, errors.New("Invalid class")
	}

	return block, nil
}

func (b *FenceWrap) Marshal() []byte {
	d, _ := xml.Marshal(b.Block)
	return append(d, '\n')
}
func (b *FenceWrap) Slice() text.Segment {
	return b.Segment
}

func ApplyWraps(source []byte, wraps []*FenceWrap) []byte {
	sort.Slice(wraps, func(i, j int) bool {
		return wraps[i].Segment.Start < wraps[j].Segment.Start
	})
	output := bytes.Buffer{}
	i := 0
	for _, w := range wraps {
		output.Write(source[i:w.Segment.Start])
		output.Write(w.Marshal())
		i = w.Segment.Stop
	}
	output.Write(source[i:])
	return output.Bytes()
}
