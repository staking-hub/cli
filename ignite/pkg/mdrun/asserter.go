package mdrun

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	cmdExec           = "exec"
	cmdExecBackground = "exec&"
	cmdWrite          = "write"
	cmdEdit           = "edit"
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

func (a *asserter) Assert(ctx context.Context, i Instruction) error {
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

	case cmdExec, cmdExecBackground:
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
				err := a.exec(ctx, cmds, cmd == cmdExecBackground)
				if err != nil {
					return ferr(fmt.Errorf("codeblock %v: %v", cmds, err))
				}
			}
		} else {
			// exec with args
			err := a.exec(ctx, s[1:], cmd == cmdExecBackground)
			if err != nil {
				return ferr(err)
			}
		}

	case cmdWrite:
		if len(s) != 2 {
			return ferr(errors.New("write requires one arg"))
		}
		filename := s[1]
		if i.CodeBlock == nil {
			return ferr(errors.New("write requires a codeblock"))
		}
		content := strings.Join(i.CodeBlock.Lines, "")
		err := os.WriteFile(filename, []byte(content), 0o644)
		if err != nil {
			return ferr(err)
		}

	case cmdEdit:
		if len(s) != 2 {
			return ferr(errors.New("edit requires one arg"))
		}
		filename := s[1]
		if i.CodeBlock == nil {
			return ferr(errors.New("edit requires a codeblock"))
		}
		_ = filename
		// TODO find how to edit a file from a snippet.

	default:
		return ferr(errors.New("unknow cmd"))
	}

	return nil
}

func (a *asserter) exec(ctx context.Context, cmds []string, async bool) error {
	log.Printf("exec(async=%t) %s", async, strings.Join(cmds, " "))
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
		a.wd = filepath.Join(a.wd, path)
		return nil
	}
	// Other than cd command
	var args []string
	if len(cmds) > 1 {
		args = cmds[1:]
	}
	if async {
		go exec.CommandContext(ctx, cmds[0], args...).Run()
		return nil
	}
	return exec.CommandContext(ctx, cmds[0], args...).Run()
}
