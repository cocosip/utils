// Package util provides a collection of utility functions.
package util

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// mockGetPath is a mock implementation of GetPathProvider for testing.
type mockGetPath struct {
	ExecutableDir string
}

func (m *mockGetPath) GetExecutableDir() string {
	return m.ExecutableDir
}

func TestGetExecutableDir(t *testing.T) {
	// 1. Test with the mock provider
	expectedDir := "/mock/path"
	PathProvider = &mockGetPath{ExecutableDir: expectedDir}

	dir := GetExecutableDir()
	if dir != expectedDir {
		t.Errorf("GetExecutableDir() with mock = %q, want %q", dir, expectedDir)
	}

	// 2. Reset to the default provider and perform a basic check
	PathProvider = &defaultGetPath{}
	dir = GetExecutableDir()
	if dir == "" {
		t.Error("GetExecutableDir() with default provider should not return an empty string")
	}
}

func TestGetCurrentDir(t *testing.T) {
	expected, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() failed: %v", err)
	}

	actual, err := GetCurrentDir()
	if err != nil {
		t.Errorf("GetCurrentDir() returned an error: %v", err)
	}

	if actual != expected {
		t.Errorf("GetCurrentDir() = %q, want %q", actual, expected)
	}
}

func TestGetFileExt(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"jpeg file", "image.jpg", ".jpg"},
		{"gzipped tarball", "archive.tar.gz", ".gz"},
		{"file with no extension", "document", ""},
		{"dotfile", ".bashrc", ""},
		{"dotfile with path", "/home/user/.vimrc", ""},
		{"file with dot in name", "my.document.pdf", ".pdf"},
		{"unicode path", "⽜⼈⽊/file.ext", ".ext"},
		{"empty path", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFileExt(tt.path)
			if got != tt.want {
				t.Errorf("GetFileExt() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetFileName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/home/user/file.txt", "file.txt"},
		{"C:\\Users\\Test\\document.docx", "document.docx"},
		{"/home/user/", "user"},
		{"file.txt", "file.txt"},
		{"/", string(filepath.Separator)},
		{"⽜⼈⽊/file.ext", "file.ext"}, // Unicode path
		{"", "."},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := GetFileName(tt.path)
			if got != tt.want {
				t.Errorf("GetFileName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsAbs(t *testing.T) {
	// Platform-independent tests
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"Relative Path", "./relative/path", false},
		{"Simple Relative", "relative", false},
		{"Unicode Relative", "⽜⼈⽊/abs/path", false},
		{"Empty Path", "", false},
	}

	// Platform-specific tests
	if runtime.GOOS == "windows" {
		tests = append(tests, []struct {
			name string
			path string
			want bool
		}{
			{"Windows Absolute", "C:\\Users\\Test", true},
			{"Windows Root", "C:\\", true},
			{"Unix-style on Windows", "/home/user", false},
		}...) // Note: The '...' is a Go syntax for appending a slice, not a string literal.
	} else {
		tests = append(tests, []struct {
			name string
			path string
			want bool
		}{
			{"Unix Absolute", "/home/user", true},
			{"Unix Root", "/", true},
		}...) // Note: The '...' is a Go syntax for appending a slice, not a string literal.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAbs(tt.path); got != tt.want {
				t.Errorf("IsAbs() = %v, want %v", got, tt.want)
			}
		})
	}
}
