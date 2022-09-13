package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func st(maxSize int64, crf string, extensions ...string) (int, error) {
	tempVideoFile := "temp.mp4"
	var allExtensions []string
	var filesChanged []string
	for _, v := range extensions {
		allExtensions = append(allExtensions, strings.ToLower(v))
		allExtensions = append(allExtensions, strings.ToUpper(v))
	}
	files, _ := getAllFilesInFolder("X:\\CopiaHD\\Videos")
	total := 0
	totalSizeBeforeMB := int64(0)
	totalSizeAfter := int64(0)
	for _, extension := range allExtensions {
		for _, f := range files {
			if strings.HasSuffix(f.info.Name(), extension) {
				sizeMB := f.info.Size() / int64(math.Pow(2, 20))
				if sizeMB > maxSize {
					fmt.Println(f.info.Name())
					totalSizeBeforeMB += sizeMB
					// ffmpeg
					cmd := exec.Command("ffmpeg", "-i", f.path, "-vcodec", "libx264", "-crf", crf, tempVideoFile)
					fmt.Println(cmd.String())
					logger.Println(cmd.String())
					cmd.Stderr = os.Stderr
					cmd.Stdout = os.Stdout
					err := cmd.Run()
					if err != nil {

						fmt.Println(err)
						logger.Println(err)
						return total, err
					}
					out, err := cmd.CombinedOutput()
					logger.Println(string(out))

					// copy all metadata from old file using exiftool
					cmd = exec.Command("exiftool", "-TagsFromFile", f.path, "-all:all>all:all", tempVideoFile)
					fmt.Println(cmd.String())
					cmd.Stderr = os.Stderr
					cmd.Stdout = os.Stdout
					err = cmd.Run()
					if err != nil {
						fmt.Println(err)
						logger.Println(err)
						return total, err
					}

					// remove old file
					err = os.Remove(f.path)
					if err != nil {
						fmt.Println(err)
						logger.Println(err)
						return total, err
					}
					directory := filepath.Dir(f.path)
					fileSplit := strings.Split(f.info.Name(), ".")
					fileNameNoExt := strings.Join(fileSplit[:len(fileSplit)-1], ".")
					destPath := filepath.Join(directory, fileNameNoExt+".mp4")
					//Rename file old to new location
					err = MoveFile(tempVideoFile, destPath)
					if err != nil {
						fmt.Println(err)
						logger.Println(err)
						return total, err
					}
					newInfo, err := os.Stat(destPath)
					if err != nil {
						fmt.Println(err)
						logger.Println(err)
						return total, err
					}
					totalSizeAfter += newInfo.Size()
					total++
					filesChanged = append(filesChanged, destPath)
					fmt.Printf("------------------> feito arquivo número: %d <----------------------\n", total)
					logger.Printf("------------------> feito arquivo número: %d <----------------------\n", total)

				}
			}
		}
	}
	fmt.Printf("Tamanho antes: %v MB\n", totalSizeBeforeMB)
	logger.Printf("Tamanho antes: %v MB\n", totalSizeBeforeMB)
	totalSizeAfterMB := totalSizeAfter / int64(math.Pow(2, 20))
	fmt.Printf("Tamanho depois: %v MB\n", totalSizeAfterMB)
	logger.Printf("Tamanho depois: %v MB\n", totalSizeAfterMB)
	fmt.Printf("Economia: %v MB\n", totalSizeBeforeMB-totalSizeAfterMB)
	logger.Printf("Economia: %v MB\n", totalSizeBeforeMB-totalSizeAfterMB)
	fmt.Printf("Lista de arquivos alterados:\n")
	logger.Printf("Lista de arquivos alterados:\n")
	for _, v := range filesChanged {
		fmt.Println(v)
		logger.Println(v)
	}
	return total, nil
}

//move file

// func runCmd(){
// 	cmd := exec.Command("youtube-dl", link, "-o", fullVideoName)

// 	//try to get all stderr line by line

// 	stderr, _ := cmd.StdoutPipe()
// 	cmd.Start()

// 	scanner := bufio.NewScanner(stderr)
// 	scanner.Split(bufio.ScanLines)
// 	var m string
// 	for scanner.Scan() {
// 		m = scanner.Text()
// 		fmt.Println(m)
// 		logger.Println(m)
// 		time.Sleep(time.Second)
// 		s := scanner.Text()
// 		fmt.Println("s externo -- >", s)
// 		for s != m {
// 			s = scanner.Text()
// 			fmt.Println("s  -- >", s)
// 			time.Sleep(time.Second * 3)
// 		}
// 	}
// 	cmd.Wait()
// }
