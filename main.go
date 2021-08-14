package main

import (
	"fmt"
	"io/ioutil"

	"github.com/metalblueberry/mdt/pkg/mdt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func main() {

	source, err := ioutil.ReadFile("README.md")
	if err != nil {
		panic(err)
	}

	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	fences, err := mdt.ParseFences(source, root)
	if err != nil {
		panic(err)
	}

	wraps, err := mdt.NewMermaidInk().WrapAll(fences)
	if err != nil {
		panic(err)
	}

	output := mdt.ApplyWraps(source, wraps)
	fmt.Print(string(output))
}
