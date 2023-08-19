package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(src string, des string) error {
	zipFile, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		path := filepath.Join(des, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err = io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
