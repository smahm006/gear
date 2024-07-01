package logger

import (
	"fmt"
	"os"
	"strings"
)

const DEFAULT_ERROR_EXIT_CODE = 1

var ErrExit = fmt.Errorf("exit")

func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

func CheckErr(err error) {
	if err == nil {
		return
	}
	switch {
	case err == ErrExit:
		fatal("", DEFAULT_ERROR_EXIT_CODE)
	default:
		msg := err.Error()
		if !strings.HasPrefix(msg, "error: ") {
			msg = fmt.Sprintf("error: %s", msg)
		}
		fatal(msg, DEFAULT_ERROR_EXIT_CODE)
	}
}
