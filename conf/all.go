package conf

func Default(path string) string {
	return joinNames(path, readDir(path, false))
}

func A(path string) string {
	names := append([]string{".", ".."}, readDir(path, true)...)
	return joinNames(path, names)
}

