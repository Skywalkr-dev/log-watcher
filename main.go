package main

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()
	watchDir := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			error := watcher.Add(path)
			if error != nil {
				return error
			}
		}
		return nil
	}
	filepath.Walk("./test-dir", watchDir)

}
