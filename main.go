package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/kjintroverted/drive-watch/drive"
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

	// handle children
	files, err := ioutil.ReadDir(os.Args[1])
	errCheck(err, "reading "+os.Args[1])

	for _, file := range files {
		if file.IsDir() {
			
			fmt.Println("Downloading from", file.Name())
			tmp, _ := ioutil.TempDir(".", file.Name())
			drive.AllDocToHTML(os.Args[1]+"/"+file.Name(), tmp)
			fmt.Println("Downloaded")
			
			fmt.Println("Converting to markdown...")
			os.Mkdir(file.Name(), 0774)
			drive.AllHTMLtoMD(tmp, file.Name())
			fmt.Println("Converted")
			os.RemoveAll(tmp)
			
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
