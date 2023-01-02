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
		// TODO add line number context
		return fmt.Errorf("assert: file '%s' cmd '%s': %v", i.Filename, i.Cmd, err)
	}
	// Set wd and restore previous at the end
	origwd, err := os.Getwd()
	if err != nil {
		return ferr(err)
	}
	if err := os.Chdir(a.wd); err != nil {
		return ferr(err)
	}
	defer os.Chdir(origwd)

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
					return ferr(fmt.Errorf("codeblock %v: %v", cmds, err))
				}
			}
		} else {
			// exec with args
			err := a.exec(s[1:])
			if err != nil {
				return ferr(err)
			}
		}

	default:
		return ferr(errors.New("unknow cmd"))
	}

	return nil
}

func (a *asserter) exec(cmds []string) error {
	if cmds[0] == "cd" {
		if len(cmds) != 2 {
			return errors.New("missing cd arg")
		}
		path := cmds[1]
		// Check path is inside a.wd
		if strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
			return fmt.Errorf("path %s must be relative w/o dots", path)
		}
		// OK, update wd
		if err := os.Chdir(path); err != nil {
			return err
		}
		return nil
	}
	// Other than cd command
	var args []string
	if len(cmds) > 1 {
		args = cmds[1:]
	}
	err := exec.Command(cmds[0], args...).Run()
	if err != nil {
		return err
	}
	return nil
}
