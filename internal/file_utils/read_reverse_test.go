// +build unit

package file_utils

import (
	"strings"
	"testing"
)

var (
	src = "Start" + "\n" +
		"Line1" + "\n" +
		"Line2" + "\n" +
		"Line3" + "\n" +
		"End"
	expected = "End" + "\n" +
		"Line3" + "\n" +
		"Line2" + "\n" +
		"Line1" + "\n" +
		"Start" + "\n"
)

func TestNewScanner(t *testing.T) {
	scanner := NewBackwardScanner(strings.NewReader(src), int64(len(src)))
	var output string
	for {
		line, _, err := scanner.Line()
		if err != nil {
			break
		}
		output = output + line + "\n"
	}
	if output != expected {
		t.Errorf("Expected, %s\n to be equal to, %s\n", output, expected)
	}
}
