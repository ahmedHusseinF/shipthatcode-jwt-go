package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
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
		//fmt.Println(line)
		parts := strings.Split(line, " ")
		rest := strings.Join(parts[1:], " ")
		if parts[0] == "ENCODE" {
			dst := base64.URLEncoding.EncodeToString([]byte(rest))
			fmt.Println(strings.TrimRight(dst, "="))
		}
		if parts[0] == "DECODE" {
			L := len(rest)
			if L%4 != 0 {
				rest += strings.Repeat("=", 4-(L%4))
			}
			src, err := base64.URLEncoding.DecodeString(rest)
			if err != nil {
				fmt.Println("Error decoding base64 string")
				panic(err)
			}
			fmt.Println(string(src))
		}
	}
}
