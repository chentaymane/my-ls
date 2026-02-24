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
