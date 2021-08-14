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

	require.Equal("```mermaid"+`
graph LR
	id
`+"```", string(source[fences[0].Segment.Start:fences[0].Segment.Stop]))

}

func TestParseFencedSequenceDiagram(t *testing.T) {
	require := require.New(t)

	source := []byte(`
# hello world
	
This is an example

` + "```mermaid" + `
sequenceDiagram
    Alice->>+John: Hello John, how are you?
    Alice->>+John: John, can you hear me?
    John-->>-Alice: Hi Alice, I can hear you!
    John-->>-Alice: I feel great!
` + "```" + `

`)

	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	fences, err := mdt.ParseFences(source, root)
	require.NoError(err)

	require.Len(fences, 1)
	require.Equal(`sequenceDiagram
    Alice->>+John: Hello John, how are you?
    Alice->>+John: John, can you hear me?
    John-->>-Alice: Hi Alice, I can hear you!
    John-->>-Alice: I feel great!
`, string(fences[0].Content()))

	require.Equal("```mermaid"+`
sequenceDiagram
    Alice->>+John: Hello John, how are you?
    Alice->>+John: John, can you hear me?
    John-->>-Alice: Hi Alice, I can hear you!
    John-->>-Alice: I feel great!
`+"```", string(source[fences[0].Segment.Start:fences[0].Segment.Stop]))

}
