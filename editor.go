// Package editor runs an external text editor to receive user input.
package editor

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

type (
	// Editor represents a command to run a text editor. This
	// command is meant to be interpreted by the shell.
	Editor string

	// File describes how to prepare a file to editing and how to
	// perform cleanup after editing.
	File interface {
		Setup() (name string, err error)
		Cleanup(name string)
	}
)

// None represents the absence of a text editor.
const None Editor = ""

// Default is an Editor returned by Resolve when $VISUAL and $EDITOR are
// unset and terminal is not dumb.
var Default Editor = "vi"

// Resolve examines environment to decide what Editor to use. If Editor
// cannot be resolved, None is returned.
func Resolve() Editor {
	var e Editor
	dumb := os.Getenv("TERM") == "dumb"
	if !dumb {
		e = Editor(os.Getenv("VISUAL"))
	}
	if e == None {
		e = Editor(os.Getenv("EDITOR"))
	}
	if e == None && !dumb {
		e = Default
	}
	return e
}

// ResolveWith calls f to resolve Editor. If f returns None, the value
// of Resolve is returned.
func ResolveWith(f func() Editor) Editor {
	if e := f(); e != None {
		return e
	}
	return Resolve()
}

// Edit runs Editor to edit a File, waits for it to complete, and returns
// resulting file content.
func (e Editor) Edit(f File) ([]byte, error) {
	if e == None {
		return nil, errors.New("editor: Editor is None")
	}

	name, err := f.Setup()
	if err != nil {
		return nil, err
	}
	defer f.Cleanup(name)

	cmd := exec.Command("/bin/sh", "-c", string(e)+` "$1"`, string(e), name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return ioutil.ReadFile(name)
}
