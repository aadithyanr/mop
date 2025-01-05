package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const Threshold = 256

func getPrunableFolders(root string) ([]string, []string) {
	entries, err := os.ReadDir(root)
	var jsPrunable []string
	var rsPrunable []string
	if err != nil {
		fmt.Println("error finding directories")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			next := filepath.Join(root, entry.Name())
			files, err := os.ReadDir(next)
			if err != nil {
				fmt.Println("Error reading directory", err)
			}
			hasNodeModules := false
			for _, file := range files {
				if file.IsDir() && file.Name() == "node_modules" {
					hasNodeModules = true
					break
				}
			}
			hasTarget := false
			for _, file := range files {
				if file.IsDir() && file.Name() == "target" {
					hasTarget = true
					break
				}
			}

			if !hasNodeModules && !hasTarget {
				continue
			}
			old := true
			for _, file := range files {
				info, err := os.Stat(filepath.Join(next, file.Name()))

				if err != nil {
					continue
				}
				if !file.IsDir() {

					if time.Now().Sub(info.ModTime()).Seconds() < Threshold*24*3600 {
						old = false
						break
					}
				}

			}
			if old {
				if hasNodeModules {
					next = filepath.Join(next, "node_modules")
					jsPrunable = append(jsPrunable, next)
				} else {
					next = filepath.Join(next, "target")
					rsPrunable = append(rsPrunable, next)
				}
			}
		}

	}
	return jsPrunable, rsPrunable
}

func main() {
	jsPrunable, rsPrunable := getPrunableFolders("../")
	fmt.Printf("\x1b[31mFound %v JavaScript projects and %v Rust projects that were last modified >%v days ago.\n\x1b[0m",
		len(jsPrunable),
		len(rsPrunable),
		Threshold)
	for _, folder := range jsPrunable {
		fmt.Print("\x1b[32m", folder, " -> Delete? (Y/N)", "\x1b[0m")
		var inp string
		fmt.Scanf("%v", inp)
	}

	for _, folder := range rsPrunable {
		fmt.Print("\x1b[33m", folder, " -> Delete? (Y/N)", "\x1b[0m")
		var inp string
		fmt.Scanf("%v", inp)
	}
}
