package editor_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/dmitrivereshchagin/editor"
)

var editFlag = flag.String("edit", "", "edit file")

type File string

func (f File) Setup() (string, error) { return string(f), nil }
func (File) Cleanup(name string)      { os.Remove(name) }

func TestResolveToDefaultIfTerminalIsNotDumb(t *testing.T) {
	defer resetenv(os.Environ())

	os.Setenv("TERM", "xterm-256color")
	os.Unsetenv("VISUAL")
	os.Unsetenv("EDITOR")

	wantEditor(t, editor.Resolve(), editor.Default)
}

func TestResolveToNoneIfTerminalIsDumb(t *testing.T) {
	defer resetenv(os.Environ())

	os.Setenv("TERM", "dumb")
	os.Unsetenv("VISUAL")
	os.Unsetenv("EDITOR")

	wantEditor(t, editor.Resolve(), editor.None)
}

func TestResolveToVisualIfTerminalIsNotDumb(t *testing.T) {
	defer resetenv(os.Environ())

	os.Setenv("TERM", "xterm-256color")
	os.Setenv("VISUAL", "vim")
	os.Setenv("EDITOR", "vim -e")

	wantEditor(t, editor.Resolve(), editor.Editor("vim"))
}

func TestResolveIgnoresVisualIfTerminalIsDumb(t *testing.T) {
	defer resetenv(os.Environ())

	os.Setenv("TERM", "dumb")
	os.Setenv("VISUAL", "emacs -nw")
	os.Unsetenv("EDITOR")

	wantEditor(t, editor.Resolve(), editor.None)
}

func TestResolveFallbacksToEditor(t *testing.T) {
	defer resetenv(os.Environ())

	os.Setenv("TERM", "xterm-256color")
	os.Unsetenv("VISUAL")
	os.Setenv("EDITOR", "code --wait")

	wantEditor(t, editor.Resolve(), editor.Editor("code --wait"))
}

func TestResolveWith(t *testing.T) {
	defer resetenv(os.Environ())

	os.Setenv("TERM", "xterm-256color")
	os.Unsetenv("VISUAL")
	os.Unsetenv("EDITOR")

	f := func() editor.Editor { return editor.Editor("nano") }
	wantEditor(t, editor.ResolveWith(f), editor.Editor("nano"))
}

func TestResolveWithFallbacksToResolve(t *testing.T) {
	defer resetenv(os.Environ())

	os.Setenv("TERM", "xterm-256color")
	os.Setenv("VISUAL", "emacs -nw")
	os.Unsetenv("EDITOR")

	f := func() editor.Editor { return editor.None }
	wantEditor(t, editor.ResolveWith(f), editor.Editor("emacs -nw"))
}

func resetenv(env []string) {
	os.Clearenv()
	for _, s := range env {
		if pair := strings.SplitN(s, "=", 2); len(pair) == 2 {
			os.Setenv(pair[0], pair[1])
		}
	}
}

func wantEditor(t *testing.T, got, want editor.Editor) {
	t.Helper()
	if got != want {
		t.Errorf("editor = %q, want %q", got, want)
	}
}

func TestEdit(t *testing.T) {
	e := editor.Editor(os.Args[0] + " -test.run TestEditor -edit")
	f := File(t.TempDir() + "/EDITME")
	content, err := e.Edit(f)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(content), "Hello, Gophers!"; got != want {
		t.Errorf("content = %q, want %q", got, want)
	}
	if _, err := os.Stat(string(f)); !os.IsNotExist(err) {
		t.Errorf("%s was not removed", f)
	}
}

func TestEditor(*testing.T) {
	if *editFlag == "" {
		return
	}
	defer os.Exit(0)
	content := []byte("Hello, Gophers!")
	if err := ioutil.WriteFile(*editFlag, content, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "TestEditor: %s\n", err)
		os.Exit(1)
	}
}
