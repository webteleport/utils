//go:build !wasi && !wasm
// +build !wasi,!wasm

package utils

import (
	"os"

	"github.com/mattn/go-isatty"
)

func Isatty() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())
}
