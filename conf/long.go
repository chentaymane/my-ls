package conf

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func modeStr(m os.FileMode) string {
	b := [10]byte{}
	switch {
	case m&os.ModeCharDevice != 0:
		b[0] = 'c'
	case m&os.ModeDevice != 0:
		b[0] = 'b'
	case m.IsDir():
		b[0] = 'd'
	case m&os.ModeSymlink != 0:
		b[0] = 'l'
	case m&os.ModeNamedPipe != 0:
		b[0] = 'p'
	case m&os.ModeSocket != 0:
		b[0] = 's'
	default:
		b[0] = '-'
	}
	p := uint32(m.Perm())
	perms := [9]struct {
		bit     uint32
		on, off byte
	}{
		{0o400, 'r', '-'},
		{0o200, 'w', '-'},
		{0o100, 'x', '-'},
		{0o040, 'r', '-'},
		{0o020, 'w', '-'},
		{0o010, 'x', '-'},
		{0o004, 'r', '-'},
		{0o002, 'w', '-'},
		{0o001, 'x', '-'},
	}
	for i, pp := range perms {
		if p&pp.bit != 0 {
			b[i+1] = pp.on
		} else {
			b[i+1] = pp.off
		}
	}
	// setuid/setgid/sticky
	if m&os.ModeSetuid != 0 {
		if b[3] == 'x' {
			b[3] = 's'
		} else {
			b[3] = 'S'
		}
	}
	if m&os.ModeSetgid != 0 {
		if b[6] == 'x' {
			b[6] = 's'
		} else {
			b[6] = 'S'
		}
	}
	if m&os.ModeSticky != 0 {
		if b[9] == 'x' {
			b[9] = 't'
		} else {
			b[9] = 'T'
		}
	}
	return string(b[:])
}

func L(path string) string {
	if path == "" {
		path = "."
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return ""
	}

	type row struct {
		mode, nlink, owner, group, size, tim, name, link string
	}
	var rows []row
	var total int64

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		fp := path + "/" + e.Name()
		info, err := os.Lstat(fp)
		if err != nil {
			continue
		}
		st := info.Sys().(*syscall.Stat_t)
		total += st.Blocks

		owner := strconv.Itoa(int(st.Uid))
		if u, err := user.LookupId(owner); err == nil {
			owner = u.Username
		}
		group := strconv.Itoa(int(st.Gid))
		if g, err := user.LookupGroupId(group); err == nil {
			group = g.Name
		}

		t := info.ModTime()
		tim := t.Format("Jan _2 15:04")
		if t.Before(time.Now().AddDate(0, -6, 0)) {
			tim = t.Format("Jan _2  2006")
		}

		var size string
		if info.Mode()&os.ModeDevice != 0 {
			maj := (st.Rdev>>8)&0xfff | (st.Rdev>>32)&^uint64(0xfff)
			min := st.Rdev&0xff | (st.Rdev>>12)&^uint64(0xff)
			size = fmt.Sprintf("%d, %3d", maj, min)
		} else {
			size = fmt.Sprintf("%d", info.Size())
		}

		lnk := ""
		if info.Mode()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(fp); err == nil {
				lnk = " -> " + target
			}
		}

		rows = append(rows, row{
			mode:  modeStr(info.Mode()),
			nlink: strconv.FormatUint(uint64(st.Nlink), 10),
			owner: owner, group: group,
			size: size, tim: tim,
			name: color(fp) + e.Name() + Reset,
			link: lnk,
		})
	}

	// column widths
	wNlink, wOwner, wGroup, wSize := 1, 1, 1, 1
	for _, r := range rows {
		if len(r.nlink) > wNlink {
			wNlink = len(r.nlink)
		}
		if len(r.owner) > wOwner {
			wOwner = len(r.owner)
		}
		if len(r.group) > wGroup {
			wGroup = len(r.group)
		}
		if len(r.size) > wSize {
			wSize = len(r.size)
		}
	}

	out := fmt.Sprintf("total %d\n", total/2)
	for _, r := range rows {
		out += fmt.Sprintf("%s %*s %-*s %-*s %*s %s %s%s\n",
			r.mode,
			wNlink, r.nlink,
			wOwner, r.owner,
			wGroup, r.group,
			wSize, r.size,
			r.tim, r.name, r.link,
		)
	}
	return out
}
