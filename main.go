package main

import (
	"fmt"
	"os"

	"main/conf"
)

func main() {
	args := os.Args[1:]

	var flags string
	var paths []string

	for _, a := range args {
		if len(a) > 1 && a[0] == '-' {
			flags += a[1:]
		} else {
			paths = append(paths, a)
		}
	}
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for i, path := range paths {
		if i > 0 {
			fmt.Println()
		}
		if len(paths) > 1 {
			fmt.Printf("%s:\n", path)
		}
		switch {
		case contains(flags, 'l'):
			fmt.Print(conf.L(path))
		case contains(flags, 'R'):
			fmt.Print(conf.R(path))
		case contains(flags, 'a'):
			fmt.Print(conf.A(path))
		case contains(flags, 't'):
			fmt.Print(conf.T(path))
		case contains(flags, 'r'):
			fmt.Print(conf.SmallR(path))
		default:
			fmt.Print(conf.Default(path))
		}
	}
}

func contains(flags string, c byte) bool {
	for i := 0; i < len(flags); i++ {
		if flags[i] == c {
			return true
		}
	}
	return false
}
