//go:build linux

package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func printLong(parent_path string, entries []os.DirEntry) {
	type row struct {
		mode    string
		nlink   string
		owner   string
		group   string
		size    string
		modTime string
		name    string
	}

	var rows []row
	var totalBlocks int64

	for _, entry := range entries {
		fullPath := parent_path + "/" + entry.Name()

		info, err := os.Lstat(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", fullPath, err)
			continue
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			continue
		}
		totalBlocks += stat.Blocks

		ownerStr := strconv.FormatUint(uint64(stat.Uid), 10)
		if u, err := user.LookupId(ownerStr); err == nil {
			ownerStr = u.Username
		}

		groupStr := strconv.FormatUint(uint64(stat.Gid), 10)
		if g, err := user.LookupGroupId(groupStr); err == nil {
			groupStr = g.Name
		}

		modTime := info.ModTime()
		var timeStr string
		if modTime.After(time.Now().AddDate(0, -6, 0)) {
			timeStr = modTime.Format("Jan _2 15:04")
		} else {
			timeStr = modTime.Format("Jan _2  2006")
		}

		name := entry.Name()
		if info.Mode()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(fullPath); err == nil {
				name += " -> " + target
			}
		}

		rows = append(rows, row{
			mode:    info.Mode().String(),
			nlink:   strconv.FormatUint(uint64(stat.Nlink), 10),
			owner:   ownerStr,
			group:   groupStr,
			size:    strconv.FormatInt(info.Size(), 10),
			modTime: timeStr,
			name:    name,
		})
	}

	fmt.Printf("total %d\n", totalBlocks/2)

	if len(rows) == 0 {
		return
	}

	maxNlink, maxOwner, maxGroup, maxSize := 1, 1, 1, 1
	for _, r := range rows {
		if len(r.nlink) > maxNlink {
			maxNlink = len(r.nlink)
		}
		if len(r.owner) > maxOwner {
			maxOwner = len(r.owner)
		}
		if len(r.group) > maxGroup {
			maxGroup = len(r.group)
		}
		if len(r.size) > maxSize {
			maxSize = len(r.size)
		}
	}

	for _, r := range rows {
		fmt.Printf("%s %*s %-*s %-*s %*s %s %s\n",
			r.mode,
			maxNlink, r.nlink,
			maxOwner, r.owner,
			maxGroup, r.group,
			maxSize, r.size,
			r.modTime,
			r.name,
		)
	}
}

// printLongFiles prints individual file paths in long format (no "total" line).
func printLongFiles(paths []string) {
	type row struct {
		mode    string
		nlink   string
		owner   string
		group   string
		size    string
		modTime string
		name    string
	}

	var rows []row

	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", path, err)
			continue
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			continue
		}

		ownerStr := strconv.FormatUint(uint64(stat.Uid), 10)
		if u, err := user.LookupId(ownerStr); err == nil {
			ownerStr = u.Username
		}

		groupStr := strconv.FormatUint(uint64(stat.Gid), 10)
		if g, err := user.LookupGroupId(groupStr); err == nil {
			groupStr = g.Name
		}

		modTime := info.ModTime()
		var timeStr string
		if modTime.After(time.Now().AddDate(0, -6, 0)) {
			timeStr = modTime.Format("Jan _2 15:04")
		} else {
			timeStr = modTime.Format("Jan _2  2006")
		}

		name := path
		if info.Mode()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(path); err == nil {
				name += " -> " + target
			}
		}

		rows = append(rows, row{
			mode:    info.Mode().String(),
			nlink:   strconv.FormatUint(uint64(stat.Nlink), 10),
			owner:   ownerStr,
			group:   groupStr,
			size:    strconv.FormatInt(info.Size(), 10),
			modTime: timeStr,
			name:    name,
		})
	}

	if len(rows) == 0 {
		return
	}

	maxNlink, maxOwner, maxGroup, maxSize := 1, 1, 1, 1
	for _, r := range rows {
		if len(r.nlink) > maxNlink {
			maxNlink = len(r.nlink)
		}
		if len(r.owner) > maxOwner {
			maxOwner = len(r.owner)
		}
		if len(r.group) > maxGroup {
			maxGroup = len(r.group)
		}
		if len(r.size) > maxSize {
			maxSize = len(r.size)
		}
	}

	for _, r := range rows {
		fmt.Printf("%s %*s %-*s %-*s %*s %s %s\n",
			r.mode,
			maxNlink, r.nlink,
			maxOwner, r.owner,
			maxGroup, r.group,
			maxSize, r.size,
			r.modTime,
			r.name,
		)
	}
}

