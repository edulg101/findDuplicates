package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
)

func x(path string) bool {
	file, err := os.Open(path)

	if err != nil {
		return false
	}
	defer file.Close()
	if strings.HasSuffix(path, ".JPG") || strings.HasSuffix(path, ".jpg") {
		_, err = jpeg.Decode(file)
		if err != nil {
			return true
		}
	}

	return false

}

func runForJpeg(rootPath string) {
	var files []string
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if x(path) && !info.IsDir() {
				files = append(files, path)
				fmt.Printf("added %v \n", path)
			}

			return nil
		})
	if err != nil {
	}
	for _, v := range files {
		fmt.Println(v)
	}
}
