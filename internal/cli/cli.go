package cli

import (
	"fmt"
	"io"
)

// Run dispatches CLI commands.
func Run(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: tg <command>")
	}
	return fmt.Errorf("command not implemented: %s", args[0])
}
