package mdt_test

import (
	"testing"

	"github.com/MetalBlueberry/mdt/pkg/mdt"
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

<div> some html block </div>

` + "```mermaid" + `
graph LR
	id
` + "```" + `

<div> some html block </div>

` + "```mermaid" + `
graph LR
	id
` + "```" + `

paragraph 
`)

	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	fences, err := mdt.ParseFences(source, root)
	require.NoError(err)

	require.Len(fences, 3)
	for _, fence := range fences {

		require.Equal(`graph LR
	id
`, string(fence.Content()))

		require.Equal("```mermaid"+`
graph LR
	id
`+"```\n", string(source[fence.Segment.Start:fence.Segment.Stop]))
	}

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
`+"```\n", string(source[fences[0].Segment.Start:fences[0].Segment.Stop]))

}

func TestParseFenced_SkipAlreadyFenced(t *testing.T) {
	require := require.New(t)

	source := []byte(`
# hello world
	
This is an example

<details class="mermaid"><summary><img src="https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBpZFxuIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0K"></img></summary><p>

` + "```mermaid" + `
graph LR
    id
` + "```" + `
</p></details>

An another paragraph here

<details class="mermaid"><summary><img src="https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBpZFxuIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0K"></img></summary><p>

` + "```mermaid" + `
graph LR
    id
` + "```" + `
</p></details>
`)

	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	fences, err := mdt.ParseFences(source, root)
	require.NoError(err)

	require.Len(fences, 0)

}

func TestParseFenced_SkipAlreadyFenced_EndFile(t *testing.T) {
	require := require.New(t)

	source := []byte(`
# hello world
	
This is an example

`)

	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	fences, err := mdt.ParseFences(source, root)
	require.NoError(err)

	require.Len(fences, 0)

}

func TestParseWrappedFences(t *testing.T) {
	require := require.New(t)

	source := []byte(`
# hello world
	
This is an example

<details class="mermaid"><summary><img src="https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBpZFxuIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0K"></img></summary><p>

` + "```mermaid" + `
graph LR
    id
` + "```" + `
</p></details>

ok
`)

	root := goldmark.DefaultParser().Parse(text.NewReader(source))
	wraps, err := mdt.ParseWrappedFences(source, root)
	require.NoError(err)

	require.Len(wraps, 1)
	for _, wrap := range wraps {

		require.Equal(`graph LR
    id
`, string(wrap.Block.Code.Mermaid))

		require.Equal(`<details class="mermaid"><summary><img src="https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBpZFxuIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0K"></img></summary><p>

`+"```mermaid"+`
graph LR
    id
`+"```"+`
</p></details>
`, string(source[wrap.Segment.Start:wrap.Segment.Stop]))

		require.Equal(`<details class="mermaid"><summary><img src="https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBpZFxuIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0K"></img></summary><p>

`+"```mermaid"+`
graph LR
    id
`+"```"+`
</p></details>
`, string(wrap.Marshal()))
	}

}
