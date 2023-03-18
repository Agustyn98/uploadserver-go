package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
)

type FileInfo struct {
	filename string
	size     int64
	lastMod  string
}

func (f FileInfo) formatSize() string {
	size := float64(f.size)
	if size <= -1 {
		return "-"
	}
	if size < 1024 {
		return fmt.Sprintf("%.0f B", size)
	}

	units := []string{"", "KB", "MB", "GB", "TB"}
	for _, unit := range units {
		if size < 1024.0 {
			return fmt.Sprintf("%.2f %s", size, unit)
		}
		size /= 1024.0
	}
	return "?"
}

func getListOfFiles(web_path string) []FileInfo {

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current dir", err)
		return nil
	}

	dir := path.Join(currentDir, web_path)

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

		var filename_dir string
		var size int64
		if entry.IsDir() {
			filename_dir = entry.Name() + "/"
			size = -1
		} else {
			filename_dir = entry.Name()
			size = info.Size()
		}

		newEntry := FileInfo{
			filename: filename_dir,
			size:     size,
			lastMod:  info.ModTime().Format("2006-01-02 15:04:05"),
		}

		if entry.IsDir() {
			files = append([]FileInfo{newEntry}, files...)
		} else {
			files = append(files, newEntry)
		}
	}
	return files
}

func getIp() string {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return ip
			}
		}
	}
	return ip
}
