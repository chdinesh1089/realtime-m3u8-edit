package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var x = 1

func getext_daterange() string {
	x += 1
	strx := strconv.Itoa(x)
	return `#EXT-X-DATERANGE:ID="` + strx + `",START-DATE="` + time.Now().Format("2006-01-02T15:04:05.000Z") + `",DURATION=5.000,X-100MSLIVE-STR="hello` + strx + `"`
}

func main() {
	input_file := os.Args[1]
	output_file := os.Args[2]

	done := make(chan bool)
	tick := time.Tick(10 * time.Second)
	<-tick

	// Process events
	go func() {
		fo, err := os.Create(output_file)
		if err != nil {
			panic(err)
		}
		stream := ""
		for {
			select {
			case <-tick:
				fo.WriteString(getext_daterange() + "\n")
			default:
				b, _ := ioutil.ReadFile(input_file)
				if stream == "" {
					stream = string(b)
					if _, err := fo.Write(b); err != nil {
						panic(err)
					}
				} else {
					if stream != string(b) {
						newS := strings.Replace(string(b), stream, "", -1)
						stream = string(b)
						fo.Write([]byte(newS))
					}
				}
			}
		}
	}()

	// Hang so program doesn't exit
	<-done
	fmt.Println("closed")
}
