package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	// "strconv"
	"strings"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func main() {

	askForDuplicate()

	// runForJpeg("E:\\Users\\Eduardo\\OneDrive\\Imagens\\Imagens no HD")

	// findAndRemoveSpecificFiles("E:\\Users\\Eduardo\\OneDrive\\Imagens\\Imagens no HD", "mp4")

	// DeletePicasaOriginals("X:\\Imagens")

	// delArray()
	// path := filepath.FromSlash("X:\\Imagens")

	// err := findAndMoveSpecificFiles(path, "X:\\Imagens\\videoTestSmall", "mov")
	// if err != nil {
	// 	panic(err)
	// }

	// getDuplicatesByYear("X:\\Imagens")

	// ScanFiles1("X:\\Imagens")

	// getDuplicates("X:\\Imagens")

	// compacta os videos
	// total, err := st(60, "22", "mov", "mp4", "mkv", "avi", "wmv")
	// fmt.Println("total:", total)
	// logger.Println("total:", total)
}

// func findAndMoveSpecificFiles(rootPath, destPath, sufix string) error {
// 	err := filepath.Walk(rootPath,
// 		func(path string, info os.FileInfo, err error) error {
// 			if err != nil {
// 				return err
// 			}
// 			if info.Size() < 4000 && info.Size() > -1 && !info.IsDir() {
// 				err := os.Remove(path)
// 				if err != nil {
// 					panic(err)
// 				} else {
// 					fmt.Printf("--> %v --- > Removed !!!!!\n", path)
// 					logger.Printf("--> %v --- > Removed !!!!!\n", path)
// 				}
// 			}

// 			if strings.HasSuffix(path, strings.ToUpper(sufix)) || strings.HasSuffix(path, strings.ToLower(sufix)) {
// 				duration, err := getDuration(path)
// 				if err != nil {
// 					duration = -1
// 				}
// 				if duration < 3 && duration >= 0 {
// 					fmt.Println(path)
// 					fmt.Println(duration)
// 					name := filepath.Base(path)
// 					destFile := filepath.Join(destPath, name)
// 					err := MoveFile(path, destFile)
// 					if err != nil {
// 						panic(err)
// 					}
// 				}

// 			}

// 			return nil
// 		})
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return nil
// }

func getDuration(path string) (float64, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	fileReader, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer fileReader.Close()

	data, err := ffprobe.ProbeReader(ctx, fileReader)
	if err != nil {
		// log.Panicf("Error getting data -> (%v) : %v", path, err)
		return -1, err
	}

	duration := data.Format.Duration()
	secs := duration.Seconds()
	return secs, nil
}

func CopyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	return nil
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func findAndRemoveSpecificFiles(rootPath, sufix string) error {
	var files []string
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, strings.ToUpper(sufix)) || strings.HasSuffix(path, strings.ToLower(sufix)) {
				duration, err := getDuration(path)
				fmt.Printf("Checking ->>>>  %v .....\n", path)
				if err != nil {
					duration = -1
				}
				if duration < 3.5 && duration >= 0 {
					err := os.Remove(path)
					if err != nil {
						panic(err)
					}
					files = append(files, path)
					fmt.Printf("file at--> %v REMOVED!!\n", path)
				}

			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}
	fmt.Println("lista dos removidos")
	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}
