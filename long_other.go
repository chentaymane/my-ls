//go:build !linux

package main

import "os"

func printLong(_ string, _ []os.DirEntry) {}

func printLongFiles(_ []string) {}
