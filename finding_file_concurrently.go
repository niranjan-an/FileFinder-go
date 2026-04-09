package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func continsText(path string, target string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, target) {
			return true
		}
	}
	return false
}

func searchConcurrent(currentPath string, target string, wg *sync.WaitGroup, limit chan struct{}, targetText string) {
	defer wg.Done()
	limit <- struct{}{}
	defer func() { <-limit }()
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
			wg.Add(1)
			go searchConcurrent(fullpath, target, wg, limit, targetText)
		} else {
			if continsText(fullpath, targetText) {
				fmt.Println("Found text in:", fullpath)
			}
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("no arguments")
	}
	var wg sync.WaitGroup
	wg.Add(1)

	root := os.Args[1]
	target := os.Args[2]
	targetText := os.Args[3]

	limit := make(chan struct{}, 50)
	go searchConcurrent(root, target, &wg, limit, targetText)
	wg.Wait()
}
