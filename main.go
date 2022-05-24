package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	input_file := os.Args[1]
	output_file := os.Args[2]
	watcher, err := fsnotify.NewWatcher()
	err = watcher.Add(input_file)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	tick := time.Tick(10 * time.Second)
	<-tick

	ext_daterange := `#EXT-X-DATERANGE:ID="999",START-DATE="2018-08-22T21:54:00.079Z",PLANNED-DURATION=30.000, SCTE35-OUT=0xFC302500000000000000FFF01405000003E77FEFFE0011FB9EFE002932E00001010100004D192A59`

	// Process events
	go func() {
		fo, err := os.Create(output_file)
		if err != nil {
			panic(err)
		}
		stream := ""
		// buf := make([]byte, 1024)
		for {
			select {
			case ev := <-watcher.Events:
				fmt.Println("\033[33m", ev, "event", "\033[0m")
				if ev.Op&fsnotify.Create == fsnotify.Create {
					b, _ := ioutil.ReadFile(input_file)
					if stream == "" {
						stream = string(b)
						if _, err := fo.Write(b); err != nil {
							panic(err)
						}
					} else {
						newS := strings.Replace(string(b), stream, "", -1)
						stream = string(b)
						fo.Write([]byte(newS))
					}
				}
			case err := <-watcher.Errors:
				panic(err)
			case <-tick:
				fo.WriteString(ext_daterange + "\n")
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
