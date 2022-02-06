package file_utils

import (
	"fmt"
	"os"
	"testing"
)

func TestNewScanner(t *testing.T) {
	f, err := os.Open("../../sample-log-files/sample.log")
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := NewScanner(f, int(fi.Size()))
	for {
		line, _, err := scanner.Line()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		fmt.Printf("%s\n", line)
	}
}
