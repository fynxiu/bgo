package gitignore

import (
	"bytes"
	"testing"
)

func TestAvailableFiles(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}
	files, err := app.availableFiles()
	if err != nil {
		t.Fatal(err)
	}
	t.Error(files)
}

func TestGenerate(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}
	var buffer bytes.Buffer
	app.generate([]string{"go", "macos"}, &buffer)
	t.Errorf(buffer.String())
}
