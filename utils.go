package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

type FileInfo struct {
	filePath string
	filename string
	size     int64
	lastMod  string
	isDir    bool
}

func getListOfFiles(web_path string) []FileInfo {

	//fmt.Printf("Getting files for dir: %s\n", web_path)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current dir", err)
		return nil
	}

	dir := path.Join(currentDir, web_path)

	//fmt.Println("Current directory:", currentDir)
	//fmt.Println("Full path:", dir)

	// ReadDir reads the named directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	files := []FileInfo{}
	// Loop over the entries
	for _, entry := range entries {
		// Get the file info using Stat()
		info, err := entry.Info()
		if err != nil {
			log.Fatal(err)
		}

		// Print the file name, date and size
		//fmt.Printf("%s - %s - %d bytes\n - %t", entry.Name(), info.ModTime(), info.Size(), info.IsDir())
		var filename_dir string
		if entry.IsDir() {
			filename_dir = entry.Name() + "/"
		} else {
			filename_dir = entry.Name()
		}
		files = append(files, FileInfo{filePath: entry.Name(),
			filename: filename_dir,
			size:     info.Size(),
			lastMod:  info.ModTime().Format("2006-01-02 15:04:05"),
			isDir:    info.IsDir(),
		})
	}
	return files
}
