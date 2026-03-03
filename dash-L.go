package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"syscall"
	"time"
)

func L(filePath string) string {
	if filePath == "" {
		filePath = "."
	}
	files, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var result string

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		fullPath := filePath + "/" + file.Name()

		// Use Lstat to NOT follow symlinks
		info, err := os.Lstat(fullPath)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			log.Fatal("Failed to get syscall.Stat_t")
		}

		// Get owner username from UID
		owner, err := user.LookupId(fmt.Sprintf("%d", stat.Uid))
		ownerName := fmt.Sprintf("%d", stat.Uid)
		if err == nil {
			ownerName = owner.Username
		}

		// Get group name from GID
		grp, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid))
		groupName := fmt.Sprintf("%d", stat.Gid)
		if err == nil {
			groupName = grp.Name
		}

		// Format time: show time if within 6 months, else show year
		modTime := info.ModTime()
		var timeStr string
		sixMonthsAgo := time.Now().AddDate(0, -6, 0)
		if modTime.After(sixMonthsAgo) {
			timeStr = modTime.Format("Jan _2 15:04")
		} else {
			timeStr = modTime.Format("Jan _2  2006")
		}

		// Handle symlink target
		symlink := ""
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(fullPath)
			if err == nil {
				symlink = " -> " + target
			}
		}

		result += fmt.Sprintf("%s %d %s %s %7d %s %s%s\n",
			info.Mode(),
			stat.Nlink,
			ownerName,
			groupName,
			info.Size(),
			timeStr,
			info.Name(),
			symlink,
		)
	}
	return result
}
