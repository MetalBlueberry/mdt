package mdt_test

import (
	"testing"

	"github.com/MetalBlueberry/mdt/pkg/mdt"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	require := require.New(t)

	wrapper := mdt.NewMermaidInk()
	wrap, err := wrapper.Wrap(&mdt.Fence{
		Mermaid: []byte(`graph LR
    id
`),
	})
	require.NoError(err)

	require.Equal(`<details class="mermaid"><summary><img src="https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBpZFxuIiwibWVybWFpZCI6IntcbiAgXCJ0aGVtZVwiOiBcImRlZmF1bHRcIlxufSIsInVwZGF0ZUVkaXRvciI6ZmFsc2UsImF1dG9TeW5jIjp0cnVlLCJ1cGRhdGVEaWFncmFtIjpmYWxzZX0K"></img></summary><p>

`+"```mermaid"+`
graph LR
    id
`+"```"+`
</p></details>
`,
		string(wrap.Marshal()))

}

func TestWrapSequenceDiagram(t *testing.T) {
	require := require.New(t)

	wrapper := mdt.NewMermaidInk()
	wrap, err := wrapper.Wrap(&mdt.Fence{
		Mermaid: []byte(`sequenceDiagram
    Alice->>+John: Hello John, how are you?
    Alice->>+John: John, can you hear me?
    John-->>-Alice: Hi Alice, I can hear you!
    John-->>-Alice: I feel great!
`),
	})
	require.NoError(err)

	require.Equal(`<details class="mermaid"><summary><img src="https://mermaid.ink/img/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgQWxpY2UtPj4rSm9objogSGVsbG8gSm9obiwgaG93IGFyZSB5b3U_XG4gICAgQWxpY2UtPj4rSm9objogSm9obiwgY2FuIHlvdSBoZWFyIG1lP1xuICAgIEpvaG4tLT4-LUFsaWNlOiBIaSBBbGljZSwgSSBjYW4gaGVhciB5b3UhXG4gICAgSm9obi0tPj4tQWxpY2U6IEkgZmVlbCBncmVhdCFcbiIsIm1lcm1haWQiOiJ7XG4gIFwidGhlbWVcIjogXCJkZWZhdWx0XCJcbn0iLCJ1cGRhdGVFZGl0b3IiOmZhbHNlLCJhdXRvU3luYyI6dHJ1ZSwidXBkYXRlRGlhZ3JhbSI6ZmFsc2V9Cg"></img></summary><p>

`+"```mermaid"+`
sequenceDiagram
    Alice->>+John: Hello John, how are you?
    Alice->>+John: John, can you hear me?
    John-->>-Alice: Hi Alice, I can hear you!
    John-->>-Alice: I feel great!
`+"```"+`
</p></details>
`,
		string(wrap.Marshal()))

}
