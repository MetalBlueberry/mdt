package mdt_test

import (
	"testing"

	"github.com/metalblueberry/mdt/pkg/mdt"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	require := require.New(t)

	wrapper := mdt.NewMermaidInk()
	wrap, err := wrapper.Wrap(&mdt.Fence{
		Code: []byte(`graph LR
    id`),
	})
	require.NoError(err)

	require.Equal(`<details class="mermaid"><summary><img src="https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBpZCIsIm1lcm1haWQiOiJ7XG4gIFwidGhlbWVcIjogXCJkZWZhdWx0XCJcbn0iLCJ1cGRhdGVFZGl0b3IiOmZhbHNlLCJhdXRvU3luYyI6dHJ1ZSwidXBkYXRlRGlhZ3JhbSI6ZmFsc2V9"></img></summary><p>`+"```mermaid"+`
graph LR
    id`+"```"+`</p></details>`,
		string(wrap.Marshal()))

}
