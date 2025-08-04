// Package util provides a collection of utility functions.
package util

import (
	"io"
	"os"
)

// FileExists checks if a file or directory exists at the given path.
// It returns true if the path exists and is accessible, false otherwise.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	// If the error is that the file does not exist, then it's a clear "false".
	// For any other error (e.g., permission denied), we cannot definitively say the file exists, so we return false.
	return !os.IsNotExist(err)
}

// CreateDir creates a directory at the specified path, including any necessary parents.
// It uses os.ModePerm for file permissions.
func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// DeleteFile removes the specified file.
// If the file does not exist, it will return an error of type *fs.PathError.
func DeleteFile(filename string) error {
	return os.Remove(filename)
}

// DeleteDir removes the specified directory and all its contents.
// If the path does not exist, it does nothing and returns nil.
func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

// ListFiles lists the names of all files (not directories) in the specified directory.
func ListFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileList := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			fileList = append(fileList, entry.Name())
		}
	}

	return fileList, nil
}

// WriteStringToFile writes the given string to the specified file.
// It creates the file if it does not exist, and truncates it if it does.
func WriteStringToFile(filename string, content string) error {
	return os.WriteFile(filename, []byte(content), 0666)
}

// ReadFile reads the entire content of a file into a byte slice.
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// CopyFile copies a file from src to dst.
// It returns the number of bytes copied and an error, if any.
func CopyFile(src, dst string) (int64, error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	return io.Copy(destFile, sourceFile)
}

// IsDir checks if the given path is a directory.
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// GetFileSize returns the size of a file in bytes.
func GetFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if fileInfo.IsDir() {
		return 0, os.ErrInvalid // Or a custom error
	}
	return fileInfo.Size(), nil
}
