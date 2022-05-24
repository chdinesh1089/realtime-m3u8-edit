package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	err = watcher.Add("stream.m3u8")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		fo, err := os.Create("streamout.m3u8")
		if err != nil {
			panic(err)
		}
		stream := ""
		// buf := make([]byte, 1024)
		for {
			select {
			case ev := <-watcher.Events:
				fmt.Println("\033[33m", ev, "event", "\033[0m")
				// if ev.Op&fsnotify.Write == fsnotify.Write {
				if ev.Op&fsnotify.Create == fsnotify.Create {
					// fi, _ := os.Open("stream.m3u8")
					b, _ := ioutil.ReadFile("stream.m3u8")
					if stream == "" {
						stream = string(b)
						if _, err := fo.Write(b); err != nil {
							panic(err)
						}
					} else {
						newS := strings.Replace(string(b), stream, "", -1)
						stream = string(b)
						fmt.Println(newS)
						fo.Write([]byte(newS))
					}
					// for {
					// 	// read a chunk
					// 	fmt.Println("hllo")
					// 	n, err := fi.Read(buf)
					// 	if err != nil && err != io.EOF {
					// 		panic(err)
					// 	}
					// 	if n == 0 {
					// 		break
					// 	}

					// 	// write a chunk
					// 	if _, err := fo.Write(buf[:n]); err != nil {
					// 		panic(err)
					// 	}
					// }
					// fi.Close()
				}
			case err := <-watcher.Errors:
				panic(err)
			}
		}
	}()

	// Hang so program doesn't exit
	<-done

	fmt.Println("closed")
	/* ... do stuff ... */
	err = watcher.Close()
	if err != nil {
		panic(err)
	}
}
