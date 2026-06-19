# my-ls

A reimplementation of the Unix `ls` command written in Go.

## Usage

```
./my-ls [flags] [file|directory ...]
```

If no path is given, the current directory is listed.

## Flags

| Flag | Long form | Description |
|------|-----------|-------------|
| `-l` | | Long format: permissions, links, owner, group, size, date, name |
| `-a` | `--all` | Include hidden entries (`.` and `..`) |
| `-R` | `--recursive` | List subdirectories recursively |
| `-r` | `--reverse` | Reverse the sort order |
| `-t` | | Sort by modification time, newest first |

Flags can be combined in any order: `-la`, `-lRt`, `-laRrt`, etc.

## Examples

```bash
./my-ls                        # list current directory
./my-ls -l                     # long format
./my-ls -la                    # long format, include hidden files
./my-ls -lR /etc               # long format, recursive
./my-ls -lt /var/log           # long format, sorted by time
./my-ls -lRr /home             # long format, recursive, reversed
./my-ls -la /dev               # list /dev with device files
./my-ls file1 dir1 dir2        # mix of files and directories
```

## Build

```bash
go build -o my-ls .
```

Requires Go 1.21+ on Linux.

## Allowed packages

`fmt` · `os` · `os/user` · `strconv` · `strings` · `syscall` · `time` · `errors` · `io/fs`

`os/exec` is **not** used anywhere in this project.

## Output

- Default: multi-column layout sized to terminal width
- `-l`: one entry per line with permissions, hard links, owner, group, size, modification time, and name
- Device files (`/dev`): size field shows `major, minor` instead of byte count
- Colors match the default GNU `ls` `LS_COLORS`:
  - **Bold blue** — directories
  - **Bold cyan** — symbolic links
  - **Bold green** — executable files
  - **Bold yellow on black** — block/character devices
  - **Yellow on black** — FIFOs / named pipes
  - **Bold magenta** — sockets
  - **White on red** — setuid files
  - **Black on yellow** — setgid files
  - **White on blue** — sticky directories
  - **Blue on green** — other-writable directories
  - **Black on green** — sticky + other-writable directories
  - **Bold red on black** — broken symbolic links
