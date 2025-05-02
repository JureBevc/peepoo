package util

import (
	"fmt"
	"os"
)

func FatalError(text string, line int, column int) {
	fmt.Println(FormatError(text, line, column))
	os.Exit(1)
}

func FormatError(text string, line int, column int) error {
	return fmt.Errorf("ln %d col %d: %s", line, column, text)
}
