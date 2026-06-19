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

// printLong prints directory contents in long (-l) format.
// It includes a "total" line showing 1 KiB blocks used.
func printLong(parent_path string, entries []os.DirEntry) {
	type row struct {
		mode    string
		nlink   string
		owner   string
		group   string
		size    string // byte count, or "major, minor" for devices
		modTime string
		name    string // colored, with symlink arrow if applicable
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

		// Size field: "major, minor" for block/char devices; byte count otherwise
		mode := info.Mode()
		var sizeStr string
		if mode&os.ModeDevice != 0 {
			maj := (stat.Rdev >> 8) & 0xfff
			min := (stat.Rdev & 0xff) | ((stat.Rdev >> 12) & 0xfff00)
			sizeStr = fmt.Sprintf("%d, %d", maj, min)
		} else {
			sizeStr = strconv.FormatInt(info.Size(), 10)
		}

		// Name: colored + symlink arrow
		color := getColor(fullPath, entry)
		name := entry.Name()
		var nameDisplay string
		if mode&os.ModeSymlink != 0 {
			target, _ := os.Readlink(fullPath)
			nameDisplay = color + name + colorReset + " -> " + target
		} else if color != "" {
			nameDisplay = color + name + colorReset
		} else {
			nameDisplay = name
		}

		rows = append(rows, row{
			mode:    info.Mode().String(),
			nlink:   strconv.FormatUint(uint64(stat.Nlink), 10),
			owner:   ownerStr,
			group:   groupStr,
			size:    sizeStr,
			modTime: timeStr,
			name:    nameDisplay,
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

		mode := info.Mode()
		var sizeStr string
		if mode&os.ModeDevice != 0 {
			maj := (stat.Rdev >> 8) & 0xfff
			min := (stat.Rdev & 0xff) | ((stat.Rdev >> 12) & 0xfff00)
			sizeStr = fmt.Sprintf("%d, %d", maj, min)
		} else {
			sizeStr = strconv.FormatInt(info.Size(), 10)
		}

		color := getColor(path, fileInfoEntry{info})

		name := path
		var nameDisplay string
		if mode&os.ModeSymlink != 0 {
			target, _ := os.Readlink(path)
			nameDisplay = color + name + colorReset + " -> " + target
		} else if color != "" {
			nameDisplay = color + name + colorReset
		} else {
			nameDisplay = name
		}

		rows = append(rows, row{
			mode:    info.Mode().String(),
			nlink:   strconv.FormatUint(uint64(stat.Nlink), 10),
			owner:   ownerStr,
			group:   groupStr,
			size:    sizeStr,
			modTime: timeStr,
			name:    nameDisplay,
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
