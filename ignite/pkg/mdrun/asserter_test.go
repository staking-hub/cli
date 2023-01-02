package mdrun_test

import (
	"context"
	_ "embed"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/mdrun"
	envtest "github.com/ignite/cli/integration"
)

var (
	//go:embed testdata/file.go
	fileGo string
	//go:embed testdata/edit.go
	editGo string
	//go:embed testdata/file_edited.go
	fileEditedGo string
)

func TestAssert(t *testing.T) {
	origwd, _ := os.Getwd()
	tests := []struct {
		name          string
		instructions  []mdrun.Instruction
		assert        func(*testing.T, mdrun.Asserter)
		expectedError string
	}{
		{
			name: "fail: empty cmd",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
			}},
			expectedError: "assert: file '01.md' cmd '': empty cmd",
		},
		{
			name: "fail: unknow content",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "xyz",
			}},
			expectedError: "assert: file '01.md' cmd 'xyz': unknow cmd",
		},
		{
			name: "fail: exec change wd without arg",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec cd",
			}},
			expectedError: "assert: file '01.md' cmd 'exec cd': missing cd arg",
		},
		{
			name: "fail: exec change wd to absolute path",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec cd /tmp",
			}},
			expectedError: "assert: file '01.md' cmd 'exec cd /tmp': path /tmp must be relative w/o dots",
		},
		{
			name: "fail: exec change wd outside initial wd #2",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec cd tmp/../..",
			}},
			expectedError: "assert: file '01.md' cmd 'exec cd tmp/../..': path tmp/../.. must be relative w/o dots",
		},
		{
			name: "ok: exec touch 1",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec touch 1",
			}},
			assert: func(t *testing.T, a mdrun.Asserter) {
				require.FileExists(t, path.Join(a.Getwd(), "1"))
			},
		},
		{
			name: "fail: exec invalid cd command",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec cd xxx",
			}},
			expectedError: "assert: file '01.md' cmd 'exec cd xxx': chdir xxx: no such file or directory",
		},
		{
			name: "fail: exec invalid command",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec foo",
			}},
			expectedError: "assert: file '01.md' cmd 'exec foo': exec: \"foo\": executable file not found in $PATH",
		},
		{
			name: "fail: single exec without code block",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec",
			}},
			expectedError: "assert: file '01.md' cmd 'exec': missing codeblock for exec",
		},
		{
			name: "ok: exec with code block",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec",
				CodeBlock: &mdrun.CodeBlock{
					Lines: []string{
						"mkdir tmp",
						"touch tmp/1",
					},
				},
			}},
			assert: func(t *testing.T, a mdrun.Asserter) {
				require.FileExists(t, path.Join(a.Getwd(), "tmp/1"))
			},
		},
		{
			name: "fail: exec with code block, invalid cd command",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec",
				CodeBlock: &mdrun.CodeBlock{
					Lines: []string{
						"cd xxx",
					},
				},
			}},
			expectedError: "assert: file '01.md' cmd 'exec': codeblock [cd xxx]: chdir xxx: no such file or directory",
		},
		{
			name: "fail: exec with code block, invalid command",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec",
				CodeBlock: &mdrun.CodeBlock{
					Lines: []string{
						"foo",
					},
				},
			}},
			expectedError: "assert: file '01.md' cmd 'exec': codeblock [foo]: exec: \"foo\": executable file not found in $PATH",
		},
		{
			name: "ok: exec with code block and $ prefix",
			instructions: []mdrun.Instruction{{
				Filename: "01.md",
				Cmd:      "exec",
				CodeBlock: &mdrun.CodeBlock{
					Lines: []string{
						"$ mkdir tmp\n",
						"$ touch tmp/1\n",
					},
				},
			}},
			assert: func(t *testing.T, a mdrun.Asserter) {
				require.FileExists(t, path.Join(a.Getwd(), "tmp/1"))
			},
		},
		{
			name: "fail: write w/o args",
			instructions: []mdrun.Instruction{
				{
					Filename: "01.md",
					Cmd:      "write",
				},
			},
			expectedError: "assert: file '01.md' cmd 'write': write requires one arg",
		},
		{
			name: "fail: write w/o codeblock",
			instructions: []mdrun.Instruction{
				{
					Filename: "01.md",
					Cmd:      "write file.go",
				},
			},
			expectedError: "assert: file '01.md' cmd 'write file.go': write requires a codeblock",
		},
		{
			name: "ok: write file",
			instructions: []mdrun.Instruction{
				{
					Cmd: "write file.go",
					CodeBlock: &mdrun.CodeBlock{
						Lines: strings.SplitAfter(fileGo, "\n"),
					},
				},
			},
			assert: func(t *testing.T, a mdrun.Asserter) {
				resFile := path.Join(a.Getwd(), "file.go")
				require.FileExists(t, resFile)
				bz, err := os.ReadFile(resFile)
				require.NoError(t, err)
				assert.Equal(t, fileGo, string(bz))
			},
		},
		{
			name: "fail: edit w/o args",
			instructions: []mdrun.Instruction{
				{
					Filename: "01.md",
					Cmd:      "edit",
				},
			},
			expectedError: "assert: file '01.md' cmd 'edit': edit requires one arg",
		},
		{
			name: "fail: edit w/o codeblock",
			instructions: []mdrun.Instruction{
				{
					Filename: "01.md",
					Cmd:      "edit file.go",
				},
			},
			expectedError: "assert: file '01.md' cmd 'edit file.go': edit requires a codeblock",
		},
		{
			name: "ok: edit file",
			instructions: []mdrun.Instruction{
				{
					Cmd: "write file.go",
					CodeBlock: &mdrun.CodeBlock{
						Lines: strings.SplitAfter(fileGo, "\n"),
					},
				},
				{
					Cmd: "edit file.go",
					CodeBlock: &mdrun.CodeBlock{
						Lines: strings.SplitAfter(editGo, "\n"),
					},
				},
			},
			assert: func(t *testing.T, a mdrun.Asserter) {
				resFile := path.Join(a.Getwd(), "file.go")
				require.FileExists(t, resFile)
				bz, err := os.ReadFile(resFile)
				require.NoError(t, err)
				assert.Equal(t, fileEditedGo, string(bz))
			},
		},
		{
			name: "ok: scaffold a chain and a module",
			instructions: []mdrun.Instruction{
				{
					Filename: "01.md",
					Cmd:      "exec",
					CodeBlock: &mdrun.CodeBlock{
						Lines: []string{
							"/tmp/ignite-tests/ignite scaffold chain hello --no-module",
						},
					},
				},
				{
					Filename: "01.md",
					Cmd:      "exec",
					CodeBlock: &mdrun.CodeBlock{
						Lines: []string{
							"cd hello",
							"/tmp/ignite-tests/ignite scaffold module mymod",
						},
					},
				},
				{
					Filename: "01.md",
					Cmd:      "exec&",
					CodeBlock: &mdrun.CodeBlock{
						Lines: []string{
							"/tmp/ignite-tests/ignite chain serve",
						},
					},
				},
			},
			assert: func(t *testing.T, a mdrun.Asserter) {
				require.FileExists(t, path.Join(a.Getwd(), "go.mod"))
				require.FileExists(t, path.Join(a.Getwd(), "/x/mymod/keeper/keeper.go"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			a, err := mdrun.DefaultAsserter()
			require.NoError(err)
			// Compile ignite binary
			envtest.New(t)
			for _, instruction := range tt.instructions {

				err := a.Assert(context.Background(), instruction)

				if tt.expectedError != "" {
					require.EqualError(err, tt.expectedError)
					return
				}
				require.NoError(err)
			}
			tt.assert(t, a)
			// Ensure the asserter doesnt alter wd after the test
			wd, _ := os.Getwd()
			assert.Equal(origwd, wd, "wd have changed")
		})
	}
}
