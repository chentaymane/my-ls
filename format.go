package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func printDirectory(parent_path string, config Config) {
	rawEntries, err := os.ReadDir(parent_path)
	if err != nil {
		panic(err)
	}

	// Filter hidden files
	var entries []os.DirEntry
	for _, entry := range rawEntries {
		if strings.HasPrefix(entry.Name(), ".") && !config.all {
			continue
		}
		entries = append(entries, entry)
	}

	// Sort: by modification time (newest first) if -t, else case-insensitive alpha
	n := len(entries)
	if config.time {
		for i := 0; i < n; i++ {
			for j := 0; j < n-i-1; j++ {
				iInfo, _ := entries[j].Info()
				jInfo, _ := entries[j+1].Info()
				if iInfo != nil && jInfo != nil && iInfo.ModTime().Before(jInfo.ModTime()) {
					entries[j], entries[j+1] = entries[j+1], entries[j]
				}
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := 0; j < n-i-1; j++ {
				if strings.ToLower(entries[j].Name()) > strings.ToLower(entries[j+1].Name()) {
					entries[j], entries[j+1] = entries[j+1], entries[j]
				}
			}
		}
	}
	if config.reverse {
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			entries[i], entries[j] = entries[j], entries[i]
		}
	}

	// Build names and dirs lists in sorted order
	var entrynames []string
	var dirs []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		}
		entrynames = append(entrynames, entry.Name())
	}

	if config.long {
		printLong(parent_path, entries)
	} else {
		printNames(entrynames)
	}

	if config.recursive {
		for _, dir := range dirs {
			fmt.Println()
			dir_path := fmt.Sprintf("%s/%s", parent_path, dir.Name())
			fmt.Println(dir_path + ":")
			printDirectory(dir_path, config)
		}
	}
}

func printNames(names []string) {
	w := getTerminalWidth()
	var nrows int = 1
	for {
		ncols := int(math.Ceil(float64(len(names)) / float64(nrows)))
		if ncols == 1 {
			for _, name := range names {
				fmt.Println(name)
			}
			return
		}

		var length int
		var offsets []int
		for i := range ncols {
			var offset int
			var colnames []string
			if (i+1)*nrows >= len(names) {
				colnames = names[i*nrows:]
			} else {
				colnames = names[i*nrows : (i+1)*nrows]
			}
			for _, name := range colnames {
				if len(name) > offset {
					offset = len(name)
				}
			}
			if length > 0 {
				length += 2
			}
			length += offset
			offsets = append(offsets, offset)
			// why counting both offsets and length ?!
		}
		if length <= int(w) {
			var out strings.Builder
			for i := range nrows {
				var lineEntities []string
				for j := range ncols {
					if j*nrows+i >= len(names) {
						break
					}
					fname := names[j*nrows+i] + strings.Repeat(" ", offsets[j]-len(names[j*nrows+i]))
					lineEntities = append(lineEntities, fname)
				}
				out.WriteString(strings.Join(lineEntities, "  ") + "\n")
			}
			fmt.Print(out.String())
			return
		}
		nrows++
	}
}
