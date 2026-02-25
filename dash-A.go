package main

import (
	"fmt"
	"log"
	"os"
)

func A(Path string) {

	files, err := os.ReadDir(Path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(".")
	fmt.Println("..")

	for _, file := range files {

		fmt.Println(file.Name())
	}
}
