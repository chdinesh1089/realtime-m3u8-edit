package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	input_file := os.Args[1]
	output_file := os.Args[2]

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
		for {
			select {
			case <-tick:
				fo.WriteString(ext_daterange + "\n")
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
