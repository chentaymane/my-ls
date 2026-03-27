package main

// bugs and missings to handle in this version of the project:
// 1. ls [FILE]
// 2. multiple flags with one dash
// 3. long listing columns alignement
// 4. flags parsing: --(double dash alone)
// 5. flags parsing: fallback mechanism (CLI flags > env vars > config file > hardcoded defaults)
// 6. should handle concurrent file read and delete ?!

type Config struct {
	long, recursive, all, reverse, time bool
}

/*
Required flags:
-l, -R, -a, -r, -t
*/

func parseFlags(args []string) ([]string, Config) {
	var paths []string // default "."
	var config Config

	for _, arg := range args {
		switch arg {
		case "-l":
			config.long = true
			// 	fmt.Print(L(path))
			//
		case "-R", "--recursive":
			config.recursive = true
			//
		case "-a", "--all":
			config.all = true
			//
		case "-r", "--reverse":
			config.reverse = true
			//
		case "-t": // different from --time ?!
			config.time = true
			//
		default:
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 {
		paths = append(paths, ".")
	}
	return paths, config
}
