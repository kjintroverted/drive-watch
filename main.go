package main

import (
	"strings"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/kjintroverted/drive-watch/drive"
)

func main() {
	source := os.Args[1]
	outSource := os.Args[2]

	watcher, err := fsnotify.NewWatcher()
	errCheck(err, "creating watcher")
	defer watcher.Close()

	done := make(chan bool)

	go handleEvents(watcher, source, outSource)

	// watch main
	err = watcher.Add(source)
	errCheck(err, "watching parent")
	fmt.Println("watching parent", source)

	// handle children
	files, err := ioutil.ReadDir(source)
	errCheck(err, "reading "+source)

	for _, file := range files {
		if file.IsDir() {
			if strings.ToLower(file.Name()) == "drafts" {
				continue
			} 
			
			fmt.Println("Downloading from", file.Name())
			tmp, _ := ioutil.TempDir(".", file.Name())
			drive.AllDocToHTML(source+"/"+file.Name(), tmp)
			fmt.Println("Downloaded")
			
			fmt.Println("Converting to markdown...")
			os.Mkdir(outSource+"/"+file.Name(), 0774)
			drive.AllHTMLtoMD(tmp, outSource+"/"+file.Name())
			fmt.Println("Converted")
			os.RemoveAll(tmp)
			
			err = watcher.Add(source + "/" + file.Name())
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

func handleEvents(watcher *fsnotify.Watcher, source, outSource string) {
	for {
		event := <-watcher.Events
		fmt.Println("EVENT:", event.Op, "for", strings.ReplaceAll(event.Name, source, ""))

		info, _ := os.Stat(event.Name)

		if info.IsDir() {
			if strings.ToLower(info.Name()) == "drafts" {
				continue
			}
			os.Mkdir(outSource+"/"+info.Name(), 0777)
			watcher.Add(event.Name)
			fmt.Println("watching", info.Name(), "directory")
			continue
		}

		paths := strings.Split(event.Name, "/")
		fileName := paths[len(paths) - 1]
		dirName := paths[len(paths) - 2]
		if strings.ToLower(dirName) == "drafts" {
			continue
		}
		
		tmp, _ := ioutil.TempDir(".", dirName)
		drive.DocToHTML(source+"/"+dirName, fileName, tmp)
		fmt.Println("Downloaded")
		
		fmt.Println("Converting to markdown...")
		os.Mkdir(outSource+"/"+dirName, 0774)
		drive.AllHTMLtoMD(tmp, outSource+"/"+dirName)
		fmt.Println("Converted")
		os.RemoveAll(tmp)
	}
}
