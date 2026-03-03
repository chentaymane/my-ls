package main

import (
	"fmt"
	"os"
	"strings"
)

func R(path string) string {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	type item struct {
		name  string
		entry os.DirEntry
	}

	var items []item
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		items = append(items, item{name, entry})
	}

	// Sort case-insensitively
	n := len(items)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if strings.ToLower(items[j].name) > strings.ToLower(items[j+1].name) {
				items[j], items[j+1] = items[j+1], items[j]
			}
		}
	}

	var names []string
	var dirs []string
	for _, it := range items {
		color := getColor(path+"/"+it.name, it.entry)
		names = append(names, fmt.Sprintf("%s%s%s", color, it.name, Reset))
		if it.entry.IsDir() {
			dirs = append(dirs, path+"/"+it.name)
		}
	}
	var result string
	fmt.Printf("%s:\n", path)
	if len(names) > 0 {
		result += strings.Join(names, "  ")
	}
	result += "\n"
	for _, dir := range dirs {
		R(dir)
	}
	return result
}
