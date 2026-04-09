package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func search(currentPath string, target string) {
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return
	}
	for _, entry := range entries {
		fullpath := filepath.Join(currentPath, entry.Name())
		if entry.Name() == target {
			fmt.Println(fullpath)
			return
		}
		if entry.IsDir() {
			search(fullpath, target)
		}
	}
}
func main() {
	if len(os.Args) < 3 {
		log.Fatal("no arguments")
	}
	root := os.Args[1]
	target := os.Args[2]
	search(root, target)
}
