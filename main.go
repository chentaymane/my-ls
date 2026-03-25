package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

/*

-l
-R
-a
-r
-t

*/

func main() {
	var dirPath string

	if len(os.Args) > 1 {
		dirPath = os.Args[1]
	} else {
		dirPath = "."
	}

	if len(os.Args) > 1 && os.Args[1] == "-l" {
		path := "."
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		fmt.Print(L(path))
		return
	} else if len(os.Args) > 1 && os.Args[1] == "-a" {
		path := "."
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		fmt.Println(A(path))
		return
	} else if len(os.Args) > 1 && os.Args[1] == "-r" {
		path := "."
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		fmt.Print(r(path))
		return
	} else if len(os.Args) > 1 && os.Args[1] == "-t" {
		path := "."
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		fmt.Print(t(path))
		return
	} else if len(os.Args) > 1 && os.Args[1] == "-R" {
		path := "."
		if len(os.Args) > 2 {
			path = os.Args[2]
		}
		fmt.Print(R(path))
		return
	}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		fmt.Println(file.Name())
	}
}
