package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DeletePicasaOriginals(rootPath string) {
	// root := flag.String("root", "Instagram", "dir of downloaded files")
	// flag.Parse()

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)

			// return err
		}
		if strings.HasSuffix(path, ".picasaoriginals") {
			fmt.Printf("Remove %q??\n", path)
			cmd := exec.Command("pause")
			cmd.Run()
			err := os.RemoveAll(path)
			fmt.Println(err)

		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
