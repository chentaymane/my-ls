package main

// bugs and missings to handle in this version of the project:
// 1. ls [FILE]
// 2. multiple flags with one dash
// 3. long listing columns alignement
// 4. flags parsing: --(double dash alone)
// 5. flags parsing: fallback mechanism (CLI flags > env vars > config file > hardcoded defaults)

type Config struct {
	long, recursive, all, reverse, time bool
}

func parseFlags(args []string) ([]string, Config) {
	var paths []string // default "."
	var flags Config

	for _, arg := range args {
		switch arg {
		case "-l":
			flags.long = true
			// 	fmt.Print(L(path))
			//
		case "-R", "--recursive":
			flags.recursive = true
			//
		case "-a", "--all":
			flags.all = true
			//
		case "-r", "--reverse":
			flags.reverse = true
			//
		case "-t": // different from --time ?!
			flags.time = true
			//
		default:
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 {
		paths = append(paths, ".")
	}
	return paths, flags
}
