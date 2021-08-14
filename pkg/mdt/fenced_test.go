package mdt_test

import (
	"testing"

	"github.com/metalblueberry/mdt/pkg/mdt"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestParseFenced(t *testing.T) {
	require := require.New(t)


	source := []byte(`
# hello world
	
This is an example

` + "```mermaid" + `
graph LR
	id
` + "```" + `

`)

	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	fences, err := mdt.ParseFences(source, root)
	require.NoError(err)

	require.Len(fences, 1)
	require.Equal(`graph LR
	id
`, string(fences[0].Content()))

}
