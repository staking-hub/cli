package mdrun_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/mdrun"
)

func TestAssert(t *testing.T) {
	origwd, _ := os.Getwd()
	tests := []struct {
		name          string
		instruction   mdrun.Instruction
		assert        func(*testing.T, mdrun.Asserter)
		expectedError string
	}{
		{
			name: "fail: empty cmd",
			instruction: mdrun.Instruction{
				Filename: "01.md",
			},
			expectedError: "assert: file '01.md' cmd '': empty cmd",
		},
		{
			name: "fail: unknow content",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "xyz",
			},
			expectedError: "assert: file '01.md' cmd 'xyz': unknow cmd",
		},
		{
			name: "fail: exec change wd without arg",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec cd",
			},
			expectedError: "assert: file '01.md' cmd 'exec cd': missing cd arg",
		},
		{
			name: "fail: exec change wd to absolute path",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec cd /tmp",
			},
			expectedError: "assert: file '01.md' cmd 'exec cd /tmp': path /tmp must be relative w/o dots",
		},
		{
			name: "fail: exec change wd outside initial wd #2",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec cd tmp/../..",
			},
			expectedError: "assert: file '01.md' cmd 'exec cd tmp/../..': path tmp/../.. must be relative w/o dots",
		},
		{
			name: "ok: exec touch 1",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec touch 1",
			},
			assert: func(t *testing.T, a mdrun.Asserter) {
				require.FileExists(t, path.Join(a.Getwd(), "1"))
			},
		},
		{
			name: "fail: exec invalid cd command",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec cd xxx",
			},
			expectedError: "assert: file '01.md' cmd 'exec cd xxx': chdir xxx: no such file or directory",
		},
		{
			name: "fail: exec invalid command",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec foo",
			},
			expectedError: "assert: file '01.md' cmd 'exec foo': exec: \"foo\": executable file not found in $PATH",
		},
		{
			name: "fail: single exec without code block",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec",
			},
			expectedError: "assert: file '01.md' cmd 'exec': missing codeblock for exec",
		},
		{
			name: "ok: exec with code block",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec",
				CodeBlock: &mdrun.CodeBlock{
					Lines: []string{
						"mkdir tmp",
						"touch tmp/1",
					},
				},
			},
			assert: func(t *testing.T, a mdrun.Asserter) {
				require.FileExists(t, path.Join(a.Getwd(), "tmp/1"))
			},
		},
		{
			name: "fail: exec with code block, invalid cd command",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec",
				CodeBlock: &mdrun.CodeBlock{
					Lines: []string{
						"cd xxx",
					},
				},
			},
			expectedError: "assert: file '01.md' cmd 'exec': codeblock [cd xxx]: chdir xxx: no such file or directory",
		},
		{
			name: "fail: exec with code block, invalid command",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec",
				CodeBlock: &mdrun.CodeBlock{
					Lines: []string{
						"foo",
					},
				},
			},
			expectedError: "assert: file '01.md' cmd 'exec': codeblock [foo]: exec: \"foo\": executable file not found in $PATH",
		},
		{
			name: "ok: exec with code block and $ prefix",
			instruction: mdrun.Instruction{
				Filename: "01.md",
				Cmd:      "exec",
				CodeBlock: &mdrun.CodeBlock{
					Lines: []string{
						"$ mkdir tmp",
						"$ touch tmp/1",
					},
				},
			},
			assert: func(t *testing.T, a mdrun.Asserter) {
				require.FileExists(t, path.Join(a.Getwd(), "tmp/1"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			a, err := mdrun.DefaultAsserter()
			require.NoError(err)

			err = a.Assert(tt.instruction)

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
			tt.assert(t, a)
			// Ensure the asserter doesnt alter wd after the test
			wd, _ := os.Getwd()
			assert.Equal(origwd, wd, "wd have changed")
		})
	}
}
