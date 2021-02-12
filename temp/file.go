// Package temp provides editor.File implementation for a temporary file.
package temp

import (
	"io/ioutil"
	"os"

	"github.com/dmitrivereshchagin/editor"
)

// File implements editor.File by creating and removing a temporary file.
type File struct {
	Dir     string
	Pattern string
	Content []byte
}

var _ editor.File = File{}

// Setup creates a new temporary file and writes its initial content.
func (f File) Setup() (name string, err error) {
	file, err := ioutil.TempFile(f.Dir, f.Pattern)
	if err != nil {
		return
	}
	_, err = file.Write(f.Content)
	if err1 := file.Close(); err1 != nil && err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return
	}
	return file.Name(), nil
}

// Cleanup removes the named file.
func (File) Cleanup(name string) { os.Remove(name) }
