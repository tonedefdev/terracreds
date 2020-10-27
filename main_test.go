package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestWritToFile(t *testing.T) {
	const fileName = "test.txt"
	const testText = "terracreds test sample text"
	filePath := t.TempDir() + "\\" + fileName

	file, err := os.Create(filePath)
	if err != nil {
		t.Errorf("Expected to be able to create a file at '%s", filePath)
	}
	defer file.Close()

	_, err = io.WriteString(file, testText)
	if err != nil {
		t.Errorf("Unable to write the test file at '%s'", filePath)
	} else {
		file.Sync()
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Errorf("Unable to read the test file at '%s'", filePath)
		} else {
			if string(data) != testText {
				t.Errorf("Expected string '%s' got '%s'", testText, data)
			}
		}
	}
}
