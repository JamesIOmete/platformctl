package output

import (
	"fmt"
	"io"
	"os"
)

// PrintLines writes each line to the writer with a newline.
func PrintLines(w io.Writer, lines ...string) {
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
}

// PrintError writes an error message to stderr.
func PrintError(lines ...string) {
	for _, line := range lines {
		fmt.Fprintln(os.Stderr, line)
	}
}
