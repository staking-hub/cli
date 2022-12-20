package mdrun

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	cmdExec = "exec"
)

func DefaultAsserter() (Asserter, error) {
	wd, err := os.MkdirTemp(os.TempDir(), "mdrun")
	if err != nil {
		return nil, fmt.Errorf("DefaultAsserter: %w", err)
	}
	return &asserter{
		wd: wd,
	}, nil
}

type asserter struct {
	wd string
}

func (a asserter) Getwd() string {
	return a.wd
}

func (a *asserter) Assert(i Instruction) error {
	ferr := func(err error) error {
		// TODO add line context
		return fmt.Errorf("assert: %v", err)
	}
	if i.Cmd == "" {
		return ferr(errors.New("empty cmd"))
	}
	s := strings.Fields(i.Cmd)
	cmd := s[0]
	switch cmd {

	case cmdExec:
		if len(s) == 1 {
			// single exec requires a code block
			if i.CodeBlock == nil {
				return ferr(errors.New("missing codeblock for exec"))
			}
			for _, line := range i.CodeBlock.Lines {
				cmds := strings.Fields(line)
				if cmds[0] == "$" {
					// skip shell prefix used to illustrate command lines
					cmds = cmds[1:]
				}
				err := a.exec(cmds)
				if err != nil {
					return ferr(err)
				}
			}
		} else {
			// exec has args
			err := a.exec(s[1:])
			if err != nil {
				return ferr(err)
			}
		}

	default:
		return ferr(fmt.Errorf("unknow cmd %q", cmd))
	}

	return nil
}

func (a *asserter) exec(cmds []string) error {
	ferr := func(err error) error {
		return fmt.Errorf("exec %v: %w", cmds, err)
	}
	if cmds[0] == "cd" {
		if len(cmds) != 2 {
			return ferr(errors.New("missing cd arg"))
		}
		if err := os.Chdir(cmds[1]); err != nil {
			return ferr(err)
		}
		return nil
	}
	var args []string
	if len(cmds) > 1 {
		args = cmds[1:]
	}
	err := exec.Command(cmds[0], args...).Run()
	if err != nil {
		return ferr(err)
	}
	return nil
}
