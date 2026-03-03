package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	Reset = "\033[0m"
	Blue  = "\033[1;34m"
	Green = "\033[1;32m"
	Cyan  = "\033[1;36m"
)

type entry struct {
	name string
	info os.DirEntry
}

func getColor(path string, file os.DirEntry) string {
	if file.Type()&os.ModeSymlink != 0 {
		return Cyan
	}
	if file.IsDir() {
		return Blue
	}
	info, err := file.Info()
	if err == nil && info.Mode()&0o111 != 0 {
		return Green
	}
	return Reset
}

func sortEntries(entries []entry) {
	n := len(entries)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if strings.ToLower(entries[j].name) > strings.ToLower(entries[j+1].name) {
				entries[j], entries[j+1] = entries[j+1], entries[j]
			}
		}
	}
}

func reverseEntries(entries []entry) {
	n := len(entries)
	for i := 0; i < n/2; i++ {
		entries[i], entries[n-1-i] = entries[n-1-i], entries[i]
	}
}

func r(path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var entries []entry
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		entries = append(entries, entry{file.Name(), file})
	}

	sortEntries(entries)
	reverseEntries(entries)
	var result string
	for _, e := range entries {
		color := getColor(path+"/"+e.name, e.info)
		result += fmt.Sprintf("%s%s%s\n", color, e.name, Reset)
	}
	return result
}
