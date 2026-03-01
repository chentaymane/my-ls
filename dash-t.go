package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	name  string
	mod   time.Time
	entry os.DirEntry
}


func sortByTime(files []FileInfo) {
	n := len(files)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if files[j].mod.Before(files[j+1].mod) {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
}

func t(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			name:  name,
			mod:   info.ModTime(),
			entry: entry,
		})
	}

	sortByTime(files)

	for _, file := range files {
		color := getColor(path+"/"+file.name, file.entry)
		fmt.Printf("%s%s%s\n", color, file.name, Reset)
	}
}
