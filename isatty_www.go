//go:build wasi || wasm
// +build wasi wasm

package utils

func Isatty() bool {
	return false
}
