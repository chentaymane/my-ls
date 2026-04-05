package conf

import (
	"os"
	"strings"
)

const (
	Reset               = "\033[0m"
	Dir                 = "\033[01;34m"
	Link                = "\033[01;36m"
	Exec                = "\033[01;32m"
	Pipe                = "\033[33m"
	Socket              = "\033[01;35m"
	BlockDev            = "\033[01;33m"
	CharDev             = "\033[01;33m"
	Setuid              = "\033[37;41m"
	Setgid              = "\033[30;43m"
	Sticky              = "\033[30;42m"
	OtherWritable       = "\033[34;42m"
	StickyOtherWritable = "\033[30;42m"
	Orphan              = "\033[01;31m"
)

func color(path string) string {
	info, err := os.Lstat(path)
	if err != nil {
		return Reset
	}
	m := info.Mode()

	if m&os.ModeSymlink != 0 {
		if _, err := os.Stat(path); err != nil {
			return Orphan
		}
		return Link
	}
	if m.IsDir() {
		p := m.Perm()
		if p&0o002 != 0 && m&os.ModeSticky != 0 {
			return StickyOtherWritable
		}
		if p&0o002 != 0 {
			return OtherWritable
		}
		if m&os.ModeSticky != 0 {
			return Sticky
		}
		return Dir
	}
	if m&os.ModeNamedPipe != 0 {
		return Pipe
	}
	if m&os.ModeSocket != 0 {
		return Socket
	}
	if m&os.ModeDevice != 0 {
		if m&os.ModeCharDevice != 0 {
			return CharDev
		}
		return BlockDev
	}
	if m&os.ModeSetuid != 0 {
		return Setuid
	}
	if m&os.ModeSetgid != 0 {
		return Setgid
	}
	if m&0o111 != 0 {
		return Exec
	}
	return Reset
}

func sortNames(names []string) {
	n := len(names)
	key := func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "."))
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if key(names[j]) > key(names[j+1]) {
				names[j], names[j+1] = names[j+1], names[j]
			}
		}
	}
}

func readDir(path string, showHidden bool) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if !showHidden && e.Name()[0] == '.' {
			continue
		}
		names = append(names, e.Name())
	}
	sortNames(names)
	return names
}

func joinNames(path string, names []string) string {
	var out string
	for i, name := range names {
		out += color(path+"/"+name) + name + Reset
		if i < len(names)-1 {
			out += "  "
		}
	}
	return out + "\n"
}

func getInfo(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}
