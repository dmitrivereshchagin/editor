package editor_test

import (
	"fmt"
	"log"
	"os"

	"github.com/dmitrivereshchagin/editor"
	"github.com/dmitrivereshchagin/editor/temp"
)

func Example() {
	e := editor.ResolveWith(editorFromEnv)
	if e == editor.None {
		log.Fatal("failed to resolve text editor")
	}
	content, err := e.Edit(temp.File{Pattern: "example.*"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", content)
}

func editorFromEnv() editor.Editor {
	return editor.Editor(os.Getenv("EXAMPLE_EDITOR"))
}
