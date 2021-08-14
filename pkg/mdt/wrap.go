package mdt

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"unicode"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type FenceWrap struct {
	Segment text.Segment
	Block   Block
}

func ParseFencesWrap(source []byte, root ast.Node) ([]*FenceWrap, error) {
	fencesWrap := []*FenceWrap{}
	err := ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch tnode := n.(type) {
		case *ast.HTMLBlock:
			if !entering {
				return ast.WalkContinue, nil
			}
			lastLine := tnode.Lines().Len() - 1
			first := tnode.Lines().At(0).Start
			last := tnode.Lines().At(lastLine).Stop

			block, err := parseHTML(source[first:last])
			if err != nil {
				return ast.WalkStop, err
			}
			fencesWrap = append(fencesWrap, &FenceWrap{
				Segment: text.NewSegment(first, last),
				Block:   block,
			})

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
			Code: fmt.Sprintf("```mermaid\n%s```", string(code)),
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
	Code string `xml:",innerxml"`
}

func parseHTML(src []byte) (Block, error) {
	block := Block{}

	decoder := xml.NewDecoder(bytes.NewReader(src))

	err := decoder.Decode(&block)
	if err != nil {
		return Block{}, err
	}

	block.Code.Code = strings.TrimRightFunc(block.Code.Code, unicode.IsSpace)

	return block, nil
}

func (b *FenceWrap) Marshal() []byte {
	d, _ := xml.Marshal(b.Block)
	return d
}
func (b *FenceWrap) Slice() text.Segment {
	return b.Segment
}