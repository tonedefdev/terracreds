package errors

import (
	"fmt"

	"github.com/fatih/color"
)

// CustomError implements a custom error interface
type CustomError struct {
	Message string
	Level   string
}

// Error returns a custom formatted error message
func (ce *CustomError) Error() string {
	return fmt.Sprintf("%s: %s\n", color.RedString(ce.Level), ce.Message)
}
