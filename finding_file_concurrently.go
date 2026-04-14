package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func containsText(path string, target string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), target) {
			return true
		}
	}
	return false
}

func searchConcurrent(ctx context.Context, cancel context.CancelFunc, currentPath string, target string, wg *sync.WaitGroup, limit chan struct{}, targetText string) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	select {
	case limit <- struct{}{}:
		defer func() { <-limit }()
	case <-ctx.Done():
		return
	}

	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fullpath := filepath.Join(currentPath, entry.Name())

		if entry.Name() == target {
			fmt.Println("Found file:", fullpath)
			cancel()
			return
		}

		if entry.IsDir() {
			wg.Add(1)
			go searchConcurrent(ctx, cancel, fullpath, target, wg, limit, targetText)
		} else {
			if containsText(fullpath, targetText) {
				fmt.Println("Found text in:", fullpath)
				cancel()
				return
			}
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Less arguments")
	}

	root := os.Args[1]
	target := os.Args[2]
	targetText := os.Args[3]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	limit := make(chan struct{}, 50)

	wg.Add(1)
	go searchConcurrent(ctx, cancel, root, target, &wg, limit, targetText)
	wg.Wait()
}
