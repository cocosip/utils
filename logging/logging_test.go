package logging

import "testing"

func TestNewFileLogger(t *testing.T) {
	_ = NewFileLogger(
		WithFilename("D:\\1.log"),
		WithMaxSize(2),
		WithMaxBackups(5),
		WithMaxAge(3),
		WithLocalTime(true),
		WithCompress(true),
		WithStdout(true),
	)

	t.Log("create new file logger success")
}
