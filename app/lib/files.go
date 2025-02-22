package lib

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// PathExists 判断文件路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// ZipDeCompress zip解压
func ZipDeCompress(zipFile, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if err := doZipDeCompress(file, destDir); err != nil {
			return fmt.Errorf("unzip Error %w", err)
		}
	}
	return nil
}

func doZipDeCompress(file *zip.File, destDir string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	if file.FileInfo().IsDir() {
		return nil
	}
	filename := filepath.Join(destDir, file.Name)
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return fmt.Errorf("mkdir ALL %s Error %s,filename:%s", filepath.Dir(filename), err.Error(), filename)
	}
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, rc)
	if err != nil {
		return err
	}
	return nil
}
