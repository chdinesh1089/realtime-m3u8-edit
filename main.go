package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func getext_daterange_func() func() string {
	var x = 0

	return func() string {
		x += 1
		strx := strconv.Itoa(x)
		return `#EXT-X-DATERANGE:ID="` + strx + `",START-DATE="` + time.Now().Format("2006-01-02T15:04:05.000Z") + `",DURATION=5.000,X-100MSLIVE-STR="hello` + strx + `"`
	}
}

func process(input_file string, output_file string, interval int) {
	fo, err := os.Create(output_file)
	var tick <-chan time.Time
	if err != nil {
		panic(err)
	}
	getext_daterange := getext_daterange_func()
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
				tick = time.Tick(time.Duration(interval) * time.Second)
			} else {
				if stream != string(b) {
					newS := strings.Replace(string(b), stream, "", -1)
					stream = string(b)
					fo.Write([]byte(newS))
				}
			}
		}
	}
}

func main() {
	if n, err := strconv.Atoi(os.Args[1]); err != nil {
		panic(err)
	} else {
		for i := 1; i <= n; i++ {
			ip := os.Args[2*i]
			op := os.Args[2*i+1]
			go process(ip, op, 60)
		}
	}

	done := make(chan bool)

	// Hang so program doesn't exit
	<-done
	fmt.Println("closed")
}
