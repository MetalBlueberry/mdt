package mdt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func NewMermaidInk() *MermaidInk {
	return &MermaidInk{
		BaseURL: "https://mermaid.ink/img/",
		Config: InkInput{
			Mermaid:       "{\n  \"theme\": \"default\"\n}",
			UpdateEditor:  false,
			AutoSync:      true,
			UpdateDiagram: false,
		},
	}
}

type MermaidInk struct {
	BaseURL string
	Config  InkInput
}

func (ink *MermaidInk) Wrap(f *Fence) (*FenceWrap, error) {
	fmt.Println(string(f.Content()))
	b64, err := ink.Encode(f.Content())
	if err != nil {
		return nil, err
	}

	block := NewBlock(f.Content(), ink.BaseURL+b64)
	return &FenceWrap{
		Segment: f.Slice(),
		Block:   block,
	}, nil

}

func (ink *MermaidInk) Encode(code []byte) (string, error) {
	ink.Config.Code = string(code)
	b, err := json.Marshal(ink.Config)
	if err != nil {
		return "", err
	}
	fmt.Println(string(b))
	e := base64.StdEncoding.EncodeToString(b)
	return e, nil
}

type InkInput struct {
	Code          string `json:"code"`
	Mermaid       string `json:"mermaid"`
	UpdateEditor  bool   `json:"updateEditor"`
	AutoSync      bool   `json:"autoSync"`
	UpdateDiagram bool   `json:"updateDiagram"`
}
