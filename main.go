package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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

func sign(header, payload, key []byte) []byte {
	hasher := hmac.New(sha256.New, key)
	str := safeEncodeString(string(header)) + "." + safeEncodeString(string(payload))
	n, err := hasher.Write([]byte(str))
	if err != nil || n != len(str) {
		panic("error writing to hasher")
	}
	return hasher.Sum(nil)
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
		parts := strings.Split(line, " ")
		cmd := parts[0]
		rest := strings.Join(parts[1:], " ")

		if cmd == "SIGN" {
			parts = strings.Split(rest, "|")
			headerB64 := safeEncodeString(parts[0])
			payloadB64 := safeEncodeString(parts[1])
			key, err := hex.DecodeString(strings.TrimSpace(parts[2]))
			if err != nil {
				panic("invalid hex key")
			}
			hasher := hmac.New(sha256.New, key)
			str := headerB64 + "." + payloadB64
			n, err := hasher.Write([]byte(str))
			if err != nil || n != len(str) {
				panic("error writing to hasher")
			}
			sig := hasher.Sum(nil)
			jwt := str + "." + safeEncodeString(string(sig))
			fmt.Println(jwt)
		}
		if cmd == "VERIFY" {
			parts1 := strings.Split(rest, "|")
			key, err := hex.DecodeString(strings.TrimSpace(parts1[1]))
			if err != nil {
				panic("invalid hex key")
			}
			jwtParts := strings.Split(parts1[0], ".")
			headerBytes, err := safeDecodeString(jwtParts[0])
			if err != nil {
				panic("invalid base64 header")
			}
			payloadBytes, err := safeDecodeString(jwtParts[1])
			if err != nil {
				panic("invalid base64 payload")
			}
			sig := jwtParts[2]
			theSig := sign(headerBytes, payloadBytes, key)

			if hmac.Equal([]byte(sig), []byte(safeEncodeString(string(theSig)))) {
				fmt.Println("OK")
			} else {
				fmt.Println("BAD")
			}
		}

	}
}
