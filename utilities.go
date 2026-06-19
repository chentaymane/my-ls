package main

type Config struct {
	long, recursive, all, reverse, time bool
}

/*
Required flags:
-l, -R, -a, -r, -t
*/

func parseFlags(args []string) ([]string, Config) {
	var paths []string
	var config Config

	for _, arg := range args {
		if len(arg) > 1 && arg[0] == '-' && arg[1] != '-' {
			// handle combined short flags: -la, -lR, -laRt, etc.
			for _, c := range arg[1:] {
				switch c {
				case 'l':
					config.long = true
				case 'R':
					config.recursive = true
				case 'a':
					config.all = true
				case 'r':
					config.reverse = true
				case 't':
					config.time = true
				}
			}
		} else if arg == "--recursive" {
			config.recursive = true
		} else if arg == "--all" {
			config.all = true
		} else if arg == "--reverse" {
			config.reverse = true
		} else {
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 { // default
		paths = append(paths, ".")
	}
	return paths, config
}


// func NameToDirEntry(dir string) os.DirEntry { // too much code !
// 	info, err := os.Stat(dir) // or os.Lstat(path)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}
// 	return fs.FileInfoToDirEntry(info)
// }
