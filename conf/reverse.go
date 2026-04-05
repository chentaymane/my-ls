package conf

func SmallR(path string) string {
	names := readDir(path, false)
	for i, j := 0, len(names)-1; i < j; i, j = i+1, j-1 {
		names[i], names[j] = names[j], names[i]
	}
	return joinNames(path, names)
}
