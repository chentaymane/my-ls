package conf

func R(path string) string {
	var out string
	first := true
	var walk func(string)
	walk = func(p string) {
		names := readDir(p, false)
		if !first {
			out += "\n"
		}
		first = false
		out += p + ":\n"
		out += joinNames(p, names)
		for _, name := range names {
			fp := p + "/" + name
			if info, err := getInfo(fp); err == nil && info.IsDir() {
				walk(fp)
			}
		}
	}
	walk(path)
	return out
}
