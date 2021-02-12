package temp_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/dmitrivereshchagin/editor/temp"
)

func TestSetup(t *testing.T) {
	f := temp.File{Dir: t.TempDir(), Content: []byte("Hello, Gophers!")}
	name, err := f.Setup()
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(content), "Hello, Gophers!"; got != want {
		t.Errorf("content = %q, want %q", got, want)
	}
}

func TestCleanup(t *testing.T) {
	f := temp.File{Dir: t.TempDir()}
	name, err := f.Setup()
	if err != nil {
		t.Fatal(err)
	}
	f.Cleanup(name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Errorf("%s was not removed", name)
	}
}
