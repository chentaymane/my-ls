package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

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
	printNames(filenames)
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
