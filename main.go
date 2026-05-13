package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func watchEvent(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				fmt.Println("Modified file : ", event.Name)
			} else if event.Has(fsnotify.Create) {
				fmt.Println("Created file : ", event.Name)
				watcher.Add(event.Name)
			} else if event.Has(fsnotify.Remove) {
				fmt.Println("Removed file : ", event.Name)
				watcher.Remove(event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println(err)
		}
	}
}

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
	go watchEvent(watcher)
	select {}
}
