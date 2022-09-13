package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func getAllFilesInFolder(rootPath string) ([]fileInfo, error) {
	var files []fileInfo
	fmt.Println("--------------------------------------------------")
	fmt.Println("Arquivos !")
	fmt.Println("--------------------------------------------------")
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			file := fileInfo{
				path: path,
				info: info,
			}
			fmt.Println(path)

			files = append(files, file)

			return nil
		})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return files, nil
}
