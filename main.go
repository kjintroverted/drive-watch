package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	errCheck(err, "creating watcher")
	defer watcher.Close()

	done := make(chan bool)

	go handleEvents(watcher)

	// watch main
	err = watcher.Add(os.Args[1])
	errCheck(err, "watching parent")
	fmt.Println("watching parent", os.Args[1])

	// watch children
	files, err := ioutil.ReadDir(os.Args[1])
	errCheck(err, "reading "+os.Args[1])

	for _, file := range files {
		if file.IsDir() {
			err = watcher.Add(os.Args[1] + "/" + file.Name())
			errCheck(err, "watching file "+file.Name())
			fmt.Println("watching", file.Name())
		}
	}

	<-done
}

func errCheck(err error, msg string) {
	if err != nil {
		fmt.Print("ERROR", msg)
	}
}

func handleEvents(watcher *fsnotify.Watcher) {
	for {
		event := <-watcher.Events
		fmt.Println("EVENT:", event)
	}
}
