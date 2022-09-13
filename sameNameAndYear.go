package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func getDuplicatesByYear(path string) {
	files, err := ScanFiles1(path)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("duplicatesyear.txt") // creates a file at current directory
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	filteredFiles := compareFilesByYear(files)
	for _, file := range filteredFiles {
		f.WriteString("\"" + file.path + "\"" + "," + "\"" + file.path1 + "\"" + "\n")
	}
}

func ScanFiles1(rootPath string) ([]fileInfo, error) {

	var files []fileInfo
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".picasa.ini") {
				os.Remove(path)

			}

			file := fileInfo{
				size: info.Size(),
				name: info.Name(),
				path: path,
				info: info,
			}
			files = append(files, file)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return files, nil
}

func getYear(info os.FileInfo) string {
	d := info.Sys().(*syscall.Win32FileAttributeData)
	creationTime := time.Unix(0, d.LastWriteTime.Nanoseconds())
	timeString := creationTime.String()
	return string(timeString[:4])
}

// func findDuplicatesByNameAndYear(files []string) []sameFile {
// 	for _, file := range files {

// 	}
// }

func compareFilesByYear(files []fileInfo) []sameFile {
	var sameFiles []sameFile

	for i, file := range files {
		for j, file1 := range files {
			if file.path == file1.path {
				break
			}

			if file.name == file1.name {
				yearFile := getYear(file.info)
				yearFile1 := getYear(file1.info)
				if yearFile == yearFile1 {
					fmt.Println("matched --> ", file.path)
					fmt.Printf("i: %v -- j:%v\n", i, j)
					sameFile := sameFile{
						name:  file.name,
						size:  file.size,
						path:  file.path,
						path1: file1.path,
						year:  yearFile,
					}
					sameFiles = append(sameFiles, sameFile)
				}
			}
		}
	}
	return sameFiles
}
