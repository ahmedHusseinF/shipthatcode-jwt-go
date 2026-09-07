package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// TODO (structure): implement per the lesson description.
func safeDecodeString(s string) ([]byte, error) {
	L := len(s)
	if L%4 != 0 {
		s += strings.Repeat("=", 4-(L%4))
	}
	return base64.URLEncoding.DecodeString(s)
}

func safeEncodeString(s string) string {
	dst := base64.URLEncoding.EncodeToString([]byte(s))
	return strings.TrimRight(dst, "=")
}

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
		// fmt.Println(line)
		parts := strings.Split(line, "|")
		header := safeEncodeString(parts[0])
		payload := safeEncodeString(parts[1])

		fmt.Println(header + "." + payload)
	}
}
