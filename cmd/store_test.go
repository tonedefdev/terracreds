package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestNewCommandActionStore(t *testing.T) {
	if runtime.GOOS != "linux" {
		store := exec.Command("terracreds", "store", "test")

		buffer := bytes.Buffer{}
		buffer.Write([]byte("{\"token\":\"test\"}"))
		store.Stdin = &buffer

		store.Stdout = os.Stdout
		store.Stderr = os.Stderr

		err := store.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
}
