package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	// "syscall"
	// "time"
)

func askForDuplicate() {

	fmt.Println("Digite a pasta para comparar: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\r\n", "", -1)

	request := newRequest(text)
	files, err := request.ScanFiles()
	if err != nil {
		fmt.Println(err)
	}
	request.log.Printf(" Encontrei %v arquivos --\n", len(files))
	fmt.Printf(" Encontrei %v arquivos --\n", len(files))
	fmt.Println("-----")
	fmt.Print("Digite 'S' se quiser comparar aquivos independente do nome ou qualquer outra tecla para incluir o nome na comparação: ")
	text, _ = reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\r\n", "", -1)
	var filteredFiles []sameFile
	if text == "S" || text == "s" {
		filteredFiles = request.compareSizeAndSum(files)
	} else {
		filteredFiles = request.compareAll(files)
	}
	f, err := os.Create("duplicates.csv") // creates a file at current directory
	if err != nil {
		fmt.Println(err)
		request.log.Println(err)
	}
	defer f.Close()
	for _, file := range filteredFiles {
		f.WriteString(fmt.Sprintf("%v,%v\n", file.path, file.path1))
		// f.WriteString("\"" + file.path + "\"" + "," + "\"" + file.path1 + "\"" + "\n")
	}
}

// func getDuplicates(path string) {
// 	files, err := ScanFiles(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	logger.Printf(" Encontrei %v arquivos --\n", len(files))
// 	fmt.Printf(" Encontrei %v arquivos --\n", len(files))
// 	f, err := os.Create("duplicates.csv") // creates a file at current directory
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer f.Close()
// 	filteredFiles := compareFiles(files)
// 	for _, file := range filteredFiles {
// 		f.WriteString("\"" + file.path + "\"" + "," + "\"" + file.path1 + "\"" + "\n")
// 	}
// }

func compareFiles(files []fileInfo) []sameFile {
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

func ScanFiles(rootPath string) ([]fileInfo, error) {
	var files []fileInfo
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
		log.Println(err)
	}
	return files, err
}

func checkSum(path string) ([16]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return [16]byte{}, err
	}
	sum := md5.Sum(data)
	return sum, nil

}
