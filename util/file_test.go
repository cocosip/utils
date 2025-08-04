// Package util provides a collection of utility functions.
package util

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// setupTestDir creates a temporary directory for testing.
func setupTestDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "util-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
	return tempDir
}

func TestFileExists(t *testing.T) {
	tempDir := setupTestDir(t)

	// 1. Test with a non-existent file
	if FileExists(filepath.Join(tempDir, "nonexistent")) {
		t.Error("FileExists should return false for a non-existent file")
	}

	// 2. Test with an existing file
	filePath := filepath.Join(tempDir, "exists.txt")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close() // Close the file handle
	if !FileExists(filePath) {
		t.Error("FileExists should return true for an existing file")
	}

	// 3. Test with an existing directory
	dirPath := filepath.Join(tempDir, "exists_dir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	if !FileExists(dirPath) {
		t.Error("FileExists should return true for an existing directory")
	}
}

func TestCreateDir(t *testing.T) {
	tempDir := setupTestDir(t)

	// 1. Create a single directory
	dirPath := filepath.Join(tempDir, "new_dir")
	if err := CreateDir(dirPath); err != nil {
		t.Errorf("CreateDir failed: %v", err)
	}
	if !FileExists(dirPath) {
		t.Error("Directory was not created")
	}

	// 2. Create nested directories
	nestedPath := filepath.Join(tempDir, "parent", "child")
	if err := CreateDir(nestedPath); err != nil {
		t.Errorf("CreateDir failed for nested path: %v", err)
	}
	if !FileExists(nestedPath) {
		t.Error("Nested directory was not created")
	}

	// 3. Create a directory that already exists (should not fail)
	if err := CreateDir(dirPath); err != nil {
		t.Errorf("CreateDir failed when creating an existing directory: %v", err)
	}
}

func TestDeleteFile(t *testing.T) {
	tempDir := setupTestDir(t)

	// 1. Delete an existing file
	filePath := filepath.Join(tempDir, "file_to_delete.txt")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close() // Close the file handle before deleting
	if err := DeleteFile(filePath); err != nil {
		t.Errorf("DeleteFile failed: %v", err)
	}
	if FileExists(filePath) {
		t.Error("File was not deleted")
	}

	// 2. Delete a non-existent file (should return an error)
	if err := DeleteFile(filepath.Join(tempDir, "nonexistent.txt")); err == nil {
		t.Error("DeleteFile should return an error for a non-existent file")
	}
}

func TestDeleteDir(t *testing.T) {
	tempDir := setupTestDir(t)

	// 1. Delete an existing directory with content
	dirPath := filepath.Join(tempDir, "dir_to_delete")
	if err := CreateDir(dirPath); err != nil {
		t.Fatalf("Failed to create test dir: %v", err)
	}
	filePath := filepath.Join(dirPath, "test.txt")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close() // Close the file handle

	if err := DeleteDir(dirPath); err != nil {
		t.Errorf("DeleteDir failed: %v", err)
	}
	if FileExists(dirPath) {
		t.Error("Directory was not deleted")
	}

	// 2. Delete a non-existent directory (should not fail)
	if err := DeleteDir(filepath.Join(tempDir, "nonexistent_dir")); err != nil {
		t.Errorf("DeleteDir should not fail for a non-existent directory, but got: %v", err)
	}
}

func TestListFiles(t *testing.T) {
	tempDir := setupTestDir(t)

	// 1. Setup directory with files and subdirectories
	f1, err := os.Create(filepath.Join(tempDir, "file1.txt"))
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	f1.Close()
	f2, err := os.Create(filepath.Join(tempDir, "file2.log"))
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	f2.Close()
	if err := os.Mkdir(filepath.Join(tempDir, "sub_dir"), 0755); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// 2. Test listing files
	files, err := ListFiles(tempDir)
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	expected := []string{"file1.txt", "file2.log"}
	sort.Strings(files) // Sort for consistent comparison
	sort.Strings(expected)

	if !reflect.DeepEqual(files, expected) {
		t.Errorf("ListFiles() got = %v, want %v", files, expected)
	}

	// 3. Test with an empty directory
	emptyDir := filepath.Join(tempDir, "empty")
	if err := os.Mkdir(emptyDir, 0755); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}
	files, err = ListFiles(emptyDir)
	if err != nil {
		t.Fatalf("ListFiles failed for empty dir: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("ListFiles() for empty dir should be empty, got %v", files)
	}

	// 4. Test with a non-existent directory
	_, err = ListFiles(filepath.Join(tempDir, "nonexistent"))
	if err == nil {
		t.Error("ListFiles should return an error for a non-existent directory")
	}
}

func TestWriteAndReadFile(t *testing.T) {
	tempDir := setupTestDir(t)
	filePath := filepath.Join(tempDir, "test_write.txt")
	content := "Hello, World!\nThis is a test."

	if err := WriteStringToFile(filePath, content); err != nil {
		t.Fatalf("WriteStringToFile failed: %v", err)
	}

	readContent, err := ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if string(readContent) != content {
		t.Errorf("Content mismatch: got %q, want %q", string(readContent), content)
	}
}

func TestCopyFile(t *testing.T) {
	tempDir := setupTestDir(t)
	srcPath := filepath.Join(tempDir, "source.txt")
	dstPath := filepath.Join(tempDir, "destination.txt")
	content := "Copy me!"

	if err := WriteStringToFile(srcPath, content); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	bytesCopied, err := CopyFile(srcPath, dstPath)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	if bytesCopied != int64(len(content)) {
		t.Errorf("Incorrect number of bytes copied: got %d, want %d", bytesCopied, len(content))
	}

	readContent, err := ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(readContent) != content {
		t.Errorf("Content mismatch after copy: got %q, want %q", string(readContent), content)
	}
}

func TestIsDir(t *testing.T) {
	tempDir := setupTestDir(t)
	filePath := filepath.Join(tempDir, "test_file.txt")

	if err := WriteStringToFile(filePath, "hello"); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	isDir, err := IsDir(tempDir)
	if err != nil {
		t.Fatalf("IsDir failed for directory: %v", err)
	}
	if !isDir {
		t.Error("IsDir should return true for a directory")
	}

	isDir, err = IsDir(filePath)
	if err != nil {
		t.Fatalf("IsDir failed for file: %v", err)
	}
	if isDir {
		t.Error("IsDir should return false for a file")
	}
}

func TestGetFileSize(t *testing.T) {
	tempDir := setupTestDir(t)
	filePath := filepath.Join(tempDir, "test_size.txt")
	content := "12345"

	if err := WriteStringToFile(filePath, content); err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	size, err := GetFileSize(filePath)
	if err != nil {
		t.Fatalf("GetFileSize failed: %v", err)
	}

	if size != int64(len(content)) {
		t.Errorf("GetFileSize incorrect: got %d, want %d", size, len(content))
	}

	// Test on a directory
	_, err = GetFileSize(tempDir)
	if err == nil {
		t.Error("GetFileSize should return an error for a directory")
	}
}
