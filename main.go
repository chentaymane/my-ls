package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

/*
Required flags:
-l, -R, -a, -r, -t
*/

func main() {
	paths, _ := parseFlags(os.Args[1:])

	if len(paths) == 1 {
		singlePathLogic(paths[0])
	} else {
		for i, path := range paths {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", path)
			singlePathLogic(path)
		}
	}
}

func singlePathLogic(dirPath string) {
	all_files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, file := range all_files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		files = append(files, file.Name())
	}
	fmt.Println(strings.Join(files, "  "))
}
