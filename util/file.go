package util

import "os"

func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func DeleteFile(filename string) error {
	return os.Remove(filename)
}

func DeleteDir(path string) error {
	if !FileExists(path) {
		return nil
	}
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}

func FileDelete(filename string) error {
	if !FileExists(filename) {
		return nil
	}
	return os.Remove(filename)
}

func ListFiles(dir string) ([]string, error) {
	fileList := make([]string, 0)
	files, err := os.ReadDir(dir)
	if err != nil {
		return fileList, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileList = append(fileList, f.Name())
	}

	return fileList, nil
}
