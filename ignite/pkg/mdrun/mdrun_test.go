package mdrun_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/mdrun"
	"github.com/ignite/cli/ignite/pkg/mdrun/mocks"
)

func TestInspect(t *testing.T) {
	tests := []struct {
		name   string
		folder string
		setup  func(*mocks.Asserter)
	}{
		{
			name:   "single file",
			folder: "single_file",
			setup: func(a *mocks.Asserter) {
				a.EXPECT().Assert(mdrun.Instruction{
					Cmd: "exec",
					CodeBlock: &mdrun.CodeBlock{
						Lang: "bash",
						Lines: []string{
							"$ ls\n",
							"$ touch 42\n",
						},
					},
				}).Return(nil)
			},
		},
		{
			name:   "multiple files",
			folder: "multiple_files",
			setup: func(a *mocks.Asserter) {
				a.EXPECT().Assert(mdrun.Instruction{
					Cmd: "exec",
					CodeBlock: &mdrun.CodeBlock{
						Lang: "bash",
						Lines: []string{
							"$ ls\n",
						},
					},
				}).Return(nil)
				a.EXPECT().Assert(mdrun.Instruction{
					Cmd: "write src/hello.go",
					CodeBlock: &mdrun.CodeBlock{
						Lang: "go",
						Lines: []string{
							"fmt.Println(\"Hello\")\n",
						},
					},
				}).Return(nil)
				a.EXPECT().Assert(mdrun.Instruction{
					Cmd: "exec ignite chain serve",
				}).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asserter := mocks.NewAsserter(t)
			tt.setup(asserter)

			err := mdrun.Inspect(path.Join("testdata", tt.folder), asserter)

			require.NoError(t, err)
		})
	}
}
