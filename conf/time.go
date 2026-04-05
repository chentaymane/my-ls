package conf

import (
	"os"
	"sort"
)

func T(path string) string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return ""
	}

	type entry struct {
		name string
		mod  int64
	}
	var list []entry
	for _, e := range entries {
		if e.Name()[0] == '.' {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		list = append(list, entry{e.Name(), info.ModTime().Unix()})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].mod > list[j].mod })

	names := make([]string, len(list))
	for i, e := range list {
		names[i] = e.name
	}
	return joinNames(path, names)
}
