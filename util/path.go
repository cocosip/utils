// Package util provides a collection of utility functions.
package util

import (
	"os"
	"path/filepath"
	"strings"
)

// GetPathProvider defines the interface for path retrieval operations.
// This interface is used to allow for mocking in unit tests.
type GetPathProvider interface {
	GetExecutableDir() string
}

// defaultGetPath is the default implementation of GetPathProvider.
type defaultGetPath struct{}

// PathProvider is the global instance of the path provider.
// You can replace this with a mock implementation for testing purposes.
var PathProvider GetPathProvider = &defaultGetPath{}

// GetExecutableDir returns the directory of the currently running executable.
// It falls back to "." if the executable path cannot be determined.
func (d *defaultGetPath) GetExecutableDir() string {
	// os.Executable() returns the path to the current executable.
	ex, err := os.Executable()
	if err != nil {
		return "."
	}
	// filepath.Dir returns the directory part of the path.
	return filepath.Dir(ex)
}

// GetExecutableDir is a convenience wrapper that uses the global PathProvider.
func GetExecutableDir() string {
	return PathProvider.GetExecutableDir()
}

// GetCurrentDir returns the current working directory.
func GetCurrentDir() (string, error) {
	return os.Getwd()
}

// GetFileExt returns the file extension of the path.
// It handles special cases like ".bashrc" (no extension) correctly.
func GetFileExt(path string) string {
	base := filepath.Base(path)
	// If the filename starts with a dot and contains no other dots, it's a dotfile, not an extension.
	if len(base) > 1 && base[0] == '.' && !strings.Contains(base[1:], ".") {
		return ""
	}
	return filepath.Ext(path)
}

// GetFileName returns the last element of a path.
// For example, GetFileName("/foo/bar/baz.txt") returns "baz.txt".
func GetFileName(path string) string {
	return filepath.Base(path)
}

// IsAbs reports whether the path is absolute.
func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}