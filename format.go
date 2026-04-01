package main
/*
import (
	"fmt"
	"math"
	"os"
	"strings"
)

func printDirectory(parent_path string, config Config) { // error ?
	entries, err := os.ReadDir(parent_path)
	if err != nil {
		panic(err)
	}

	var entrynames []string
	var dirs []os.DirEntry
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") && !config.all {
			continue
		}
		if entry.IsDir() { // unuseful if not recursive!
			dirs = append(dirs, entry)
		}
		entrynames = append(entrynames, entry.Name())
	}
	printNames(entrynames)

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

		for i := 0; i < ncols; i++ {
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
		}

		if length <= int(w) {
			var out strings.Builder

			for i := 0; i < nrows; i++ {
				var lineEntities []string

				for j := 0; j < ncols; j++ {
					if j*nrows+i >= len(names) {
						break
					}
					fname := names[j*nrows+i] +
						strings.Repeat(" ", offsets[j]-len(names[j*nrows+i]))

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
*/