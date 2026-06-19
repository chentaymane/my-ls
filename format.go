package main

import (
	"fmt"
	"os"
	"strings"
)

// ANSI color codes matching default GNU ls LS_COLORS.
const (
	colorReset = "\033[0m"

	// file types (foreground only)
	lsDir      = "\033[01;34m"   // di: bold blue
	lsSymlink  = "\033[01;36m"   // ln: bold cyan
	lsExec     = "\033[01;32m"   // ex: bold green
	lsSocket   = "\033[01;35m"   // so: bold magenta
	lsBroken   = "\033[40;31;01m" // or: bold red on black (broken symlink)

	// file types with background
	lsBlockDev  = "\033[40;33;01m" // bd: bold yellow on black
	lsCharDev   = "\033[40;33;01m" // cd: bold yellow on black
	lsPipe      = "\033[40;33m"    // pi: yellow on black
	lsSetuid    = "\033[37;41m"    // su: white on red
	lsSetgid    = "\033[30;43m"    // sg: black on yellow
	lsSticky    = "\033[37;44m"    // st: white on blue (sticky, not OW)
	lsOtherWr   = "\033[34;42m"    // ow: blue on green (other-writable)
	lsStickyOW  = "\033[30;42m"    // tw: black on green (sticky + OW)
)

// dotEntry is a synthetic DirEntry for "." and ".." produced by the -a flag.
type dotEntry struct {
	n    string // "." or ".."
	path string // real filesystem path, used by Info()
}

func (d dotEntry) Name() string               { return d.n }
func (d dotEntry) IsDir() bool                { return true }
func (d dotEntry) Type() os.FileMode          { return os.ModeDir }
func (d dotEntry) Info() (os.FileInfo, error) { return os.Lstat(d.path) }

// fileInfoEntry wraps os.FileInfo to implement os.DirEntry.
// Used so printLongFiles can call getColor without a real DirEntry.
type fileInfoEntry struct{ fi os.FileInfo }

func (f fileInfoEntry) Name() string               { return f.fi.Name() }
func (f fileInfoEntry) IsDir() bool                { return f.fi.IsDir() }
func (f fileInfoEntry) Type() os.FileMode          { return f.fi.Mode().Type() }
func (f fileInfoEntry) Info() (os.FileInfo, error) { return f.fi, nil }

// getColor returns the ANSI prefix for an entry, or "" for plain files.
// Matches the default GNU ls LS_COLORS (di, ln, ex, bd, cd, pi, so, su, sg, tw, ow, st, or).
func getColor(fullPath string, entry os.DirEntry) string {
	t := entry.Type()

	// Symlinks
	if t&os.ModeSymlink != 0 {
		if _, err := os.Stat(fullPath); err != nil {
			return lsBroken
		}
		return lsSymlink
	}

	// Directories: check sticky + other-writable bits
	if t.IsDir() {
		info, _ := entry.Info()
		if info != nil {
			mode := info.Mode()
			sticky := mode&os.ModeSticky != 0
			ow := mode&0002 != 0
			switch {
			case sticky && ow:
				return lsStickyOW
			case ow:
				return lsOtherWr
			case sticky:
				return lsSticky
			}
		}
		return lsDir
	}

	// Device files
	if t&os.ModeDevice != 0 {
		if t&os.ModeCharDevice != 0 {
			return lsCharDev
		}
		return lsBlockDev
	}

	// Named pipe / FIFO
	if t&os.ModeNamedPipe != 0 {
		return lsPipe
	}

	// Socket
	if t&os.ModeSocket != 0 {
		return lsSocket
	}

	// Regular files: check setuid, setgid, executable
	info, err := entry.Info()
	if err == nil && info != nil {
		mode := info.Mode()
		switch {
		case mode&os.ModeSetuid != 0:
			return lsSetuid
		case mode&os.ModeSetgid != 0:
			return lsSetgid
		case mode&0111 != 0:
			return lsExec
		}
	}

	return ""
}

func printDirectory(parent_path string, config Config) {
	rawEntries, err := os.ReadDir(parent_path)
	if err != nil {
		panic(err)
	}

	// When -a: prepend synthetic "." and ".." so they sort naturally
	var allRaw []os.DirEntry
	if config.all {
		allRaw = append(allRaw,
			dotEntry{".", parent_path},
			dotEntry{"..", parent_path + "/.."},
		)
	}
	allRaw = append(allRaw, rawEntries...)

	// Filter hidden files (unless -a)
	var entries []os.DirEntry
	for _, entry := range allRaw {
		if strings.HasPrefix(entry.Name(), ".") && !config.all {
			continue
		}
		entries = append(entries, entry)
	}

	// Sort: newest-first by mod time if -t, else case-insensitive alpha
	n := len(entries)
	if config.time {
		for i := 0; i < n; i++ {
			for j := 0; j < n-i-1; j++ {
				iInfo, _ := entries[j].Info()
				jInfo, _ := entries[j+1].Info()
				if iInfo != nil && jInfo != nil && iInfo.ModTime().Before(jInfo.ModTime()) {
					entries[j], entries[j+1] = entries[j+1], entries[j]
				}
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := 0; j < n-i-1; j++ {
				if strings.ToLower(entries[j].Name()) > strings.ToLower(entries[j+1].Name()) {
					entries[j], entries[j+1] = entries[j+1], entries[j]
				}
			}
		}
	}
	if config.reverse {
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			entries[i], entries[j] = entries[j], entries[i]
		}
	}

	// Build names, colored display names, and dirs (sorted order)
	var entrynames []string
	var displaynames []string
	var dirs []os.DirEntry
	for _, entry := range entries {
		name := entry.Name()
		fullPath := parent_path + "/" + name
		color := getColor(fullPath, entry)

		entrynames = append(entrynames, name)
		if color != "" {
			displaynames = append(displaynames, color+name+colorReset)
		} else {
			displaynames = append(displaynames, name)
		}

		// Don't recurse into "." or ".." to avoid infinite loops
		if entry.IsDir() && name != "." && name != ".." {
			dirs = append(dirs, entry)
		}
	}

	if config.long {
		printLong(parent_path, entries)
	} else {
		printNames(entrynames, displaynames)
	}

	if config.recursive {
		for _, dir := range dirs {
			fmt.Println()
			dir_path := fmt.Sprintf("%s/%s", parent_path, dir.Name())
			fmt.Println(dir_path + ":")
			printDirectory(dir_path, config)
		}
	}
}

// printNames prints names in columns.
// names is used for width calculation; display (same length) is what gets printed.
func printNames(names []string, display []string) {
	if len(names) == 0 {
		return
	}
	w := getTerminalWidth()
	var nrows int = 1
	for {
		ncols := (len(names) + nrows - 1) / nrows
		if ncols == 1 {
			for i := range names {
				fmt.Println(display[i])
			}
			return
		}

		var length int
		var offsets []int
		for i := range ncols {
			var offset int
			var colnames []string
			if (i+1)*nrows >= len(names) {
				colnames = names[i*nrows:]
			} else {
				colnames = names[i*nrows : (i+1)*nrows]
			}
			for _, name := range colnames {
				if len(name) > offset {
					offset = len(name)
				}
			}
			if length > 0 {
				length += 2
			}
			length += offset
			offsets = append(offsets, offset)
		}
		if length <= int(w) {
			var out strings.Builder
			for i := range nrows {
				var lineEntities []string
				for j := range ncols {
					idx := j*nrows + i
					if idx >= len(names) {
						break
					}
					// Width based on raw name; display the colored version
					padding := strings.Repeat(" ", offsets[j]-len(names[idx]))
					lineEntities = append(lineEntities, display[idx]+padding)
				}
				out.WriteString(strings.Join(lineEntities, "  ") + "\n")
			}
			fmt.Print(out.String())
			return
		}
		nrows++
	}
}
