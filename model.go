package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type request struct {
	path string
	log  *log.Logger
}

type fileInfo struct {
	name string
	size int64
	sum  [16]byte
	path string
	year string
	info os.FileInfo
}
type sameFile struct {
	name  string
	size  int64
	sum   [16]byte
	year  string
	path  string
	path1 string
}

func (info fileInfo) getCheckSum() [16]byte {
	sum, err := checkSum(info.path)
	if err != nil {
		sum = [16]byte{}
	}
	return sum
}

func newRequest(path string) request {
	f, err := os.OpenFile("log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	logger = log.New(f, "--> ", log.LstdFlags)

	return request{
		path: path,
		log:  logger,
	}
}

func (request request) getDuplicates() {
	path := request.path
	fmt.Printf("%q\n", path)
	logger := request.log
	files, err := request.ScanFiles(path)
	if err != nil {
		panic(err)
	}
	logger.Printf(" Encontrei %v arquivos --\n", len(files))
	fmt.Printf(" Encontrei %v arquivos --\n", len(files))
	f, err := os.Create("duplicates.csv") // creates a file at current directory
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	filteredFiles := request.compareFiles(files)
	for _, file := range filteredFiles {
		f.WriteString(fmt.Sprintf("%v,%v\n", file.path, file.path1))
		// f.WriteString("\"" + file.path + "\"" + "," + "\"" + file.path1 + "\"" + "\n")
	}
}

func (request request) compareFiles(files []fileInfo) []sameFile {
	var sameFiles []sameFile

	for _, file := range files {
		for _, file1 := range files {
			if file.path == file1.path {
				break
			}

			if file.name[:5] == file1.name[:5] {
				if file.size == file1.size {
					fileCheckSum := file.getCheckSum()
					if fileCheckSum == file1.getCheckSum() && len(fileCheckSum) > 0 {
						fmt.Printf("matched-> %v  <--> %v \n", file.path, file1.path)
						sameFile := sameFile{
							name:  file.name,
							size:  file.size,
							sum:   fileCheckSum,
							path:  file.path,
							path1: file1.path,
						}
						sameFiles = append(sameFiles, sameFile)
					}
				}
			}
		}
	}
	return sameFiles
}
func (request request) ScanFiles(rootPath string) ([]fileInfo, error) {
	var files []fileInfo
	logger := request.log
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// d := info.Sys().(*syscall.Win32FileAttributeData)
			// creationTime := time.Unix(0, d.LastWriteTime.Nanoseconds())
			// timeString := creationTime.String()
			// fmt.Println(timeString)

			file := fileInfo{
				size: info.Size(),
				name: info.Name(),
				info: info,
				path: path,
			}
			if !info.IsDir() {
				files = append(files, file)

			}
			return nil
		})
	if err != nil {
		logger.Println(err)
	}
	return files, err
}
