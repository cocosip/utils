package util

import (
	"os"
	"path/filepath"
)

type GetPathInterface interface {
	GetExecutePath() string
}

type DefaultGetPath struct{}

var GetPath GetPathInterface = &DefaultGetPath{}

func (d *DefaultGetPath) GetExecutePath() string {
	if ex, err := os.Executable(); err == nil {
		return filepath.Dir(ex)
	}
	return "."
}

func GetExecutePath() string {
	return GetPath.GetExecutePath()
}
