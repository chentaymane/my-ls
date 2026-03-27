package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	paths, config := parseFlags(os.Args[1:])

	// Order of listing:
	// 1. errors
	// 2. existing files (ascending alphabetical order)
	// 3. existing dirs (asc alpha order for dirs and their contents)
	var files, dirs []string
	for _, path := range paths {
		info, err := os.Stat(path)
		// FileInfo vs DirEntry ?!
		if err != nil {
			fmt.Println(err)
			// permission denied error should be inline, and not before!?
			// or when opening folders ?!
			continue
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		} else {
			// It is likely a file, but not always!
			files = append(files, path)
		}
	}

	// when using folder title:
	// 1. -R flag is setted
	// 2. files is not empty (len(files) > 0)
	// 3. multiple folders (len(dirs) > 1)

	// Note: paths are printed just like they were passed!

	// 1. files logic: {}
	// **********************************
	if len(files) > 0 {
		fmt.Println(strings.Join(files, "  "))
	}
	// **********************************
	if config.recursive || len(files) > 0 || len(dirs) > 1 {
		for i, dir := range dirs {
			if i > 0 || len(files) > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", dir)
			singlePathLogic(dir)
		}
	} else {
		singlePathLogic(dirs[0])
	}

	// if len(files) == 0

	// 2. folders logic:

	// if len(paths) == 1 {
	// 	singlePathLogic(paths[0])
	// } else {
	// }
}

func singlePathLogic(dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	var filenames []string
	for _, file := range entries {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		filenames = append(filenames, file.Name())
	}
	fmt.Println(strings.Join(filenames, "  "))
	// actually organized in columns with 2 space between !
}
