package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO (structure): implement per the lesson description.

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		if err := sc.Err(); err != nil {
			break
		}
		line := sc.Text()
		if line == "" {
			continue
		}
		fmt.Println("TODO")

	}
}
