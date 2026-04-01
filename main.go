package main

import (
	"fmt"
	"os"

	"main/conf"
)

type Config struct {
	all       bool
	long      bool
	recursive bool
	timeSort  bool
	reverse   bool
}


func simpleList(path string) string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Sprintln(err)
	}

	var result string

	for _, entry := range entries {
		name := entry.Name()

		if len(name) > 0 && name[0] == '.' {
			continue
		}

		result += name + "  "
	}

	result += "\n"
	return result
}

//  SELECT FUNCTION
func render(path string, config Config) string {
	if config.long {
		return conf.L(path)
	}

	if config.timeSort {
		return conf.T(path)
	}

	if config.all {
		return conf.A(path)
	}

	if config.recursive {
		return conf.R(path)
	}

	return simpleList(path)
}

//  PRINT
func printPaths(paths []string, config Config) {
	for i, path := range paths {

		info, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if i > 0 {
			fmt.Println()
		}

		// show title only when needed
		if (len(paths) > 1 || config.recursive) && info.IsDir() {
			fmt.Printf("%s:\n", path)
		}

		fmt.Print(render(path, config))
	}
}

//  PARSE FLAGS
func parseFlags(args []string) ([]string, Config) {
	var config Config
	var paths []string

	for _, arg := range args {

		if len(arg) > 1 && arg[0] == '-' {
			for _, ch := range arg[1:] {
				switch ch {
				case 'l':
					config.long = true
				case 'R':
					config.recursive = true
				case 'a':
					config.all = true
				case 'r':
					config.reverse = true
				case 't':
					config.timeSort = true
				default:
					fmt.Printf("ls: invalid option -- '%c'\n", ch)
				}
			}
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 {
		paths = []string{"."}
	}

	return paths, config
}

func main() {
	paths, config := parseFlags(os.Args[1:])
	printPaths(paths, config)
}