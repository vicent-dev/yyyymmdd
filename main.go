package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {

	var dir string

	if len(os.Args) != 2 {
		dir, _ = os.Getwd()
	} else {
		dir = os.Args[1]
	}

	log.Printf("Organising %s", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(files))

	for _, f := range files {
		go func(wg *sync.WaitGroup, dir string, file os.DirEntry) {
			defer wg.Done()
			moveFile(dir, file)
		}(&wg, dir, f)
	}

	wg.Wait()

	fmt.Println("All files moved")
	os.Exit(1)
}

func moveFile(dir string, file os.DirEntry) {
	if file.IsDir() {
		return
	}

	fullPath, err := filepath.Abs(filepath.Join(dir, file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(file.Name(), ".") {
		return
	}

	fileInfo, err := file.Info()

	if err != nil {
		log.Fatal(err)
	}

	yyyymmdd := fileInfo.ModTime().Format("2006-01-02")

	// new path in "yyyy-mm-dd" folder
	newPath := filepath.Join(dir, yyyymmdd)
	newFullPath := filepath.Join(newPath, file.Name())

	err = os.MkdirAll(newPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(fullPath, newFullPath)
	if err != nil {
		log.Fatal(err)
	}
}
