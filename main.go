package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var rsaPublicKeys []string
var hs256Keys [][]byte

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

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

func signHs(header, payload, key []byte) []byte {
	hasher := hmac.New(sha256.New, key)
	str := safeEncodeString(string(header)) + "." + safeEncodeString(string(payload))
	n, err := hasher.Write([]byte(str))
	if err != nil || n != len(str) {
		panic("error writing to hasher")
	}
	return hasher.Sum(nil)
}

func signRs(header, payload, key []byte) []byte {
	return nil
}

func init() {
	rsaPublicKeys = make([]string, 0)
	hs256Keys = make([][]byte, 0)
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
		if cmd == "KEY" {
			split := strings.Split(rest, " ")
			kind := strings.TrimSpace(split[0])
			key := strings.TrimSpace(split[1])
			if kind == "oct" {
				keyBytes, err := hex.DecodeString(key)
				if err != nil {
					panic("invalid hex key")
				}
				hs256Keys = append(hs256Keys, keyBytes)
			} else if kind == "pub" {
				rsaPublicKeys = append(rsaPublicKeys, key)
			}
			fmt.Println("OK")
		}
		if cmd == "VERIFY" {
			jwtParts := strings.Split(rest, ".")
			headerBytes, err := safeDecodeString(jwtParts[0])
			if err != nil {
				panic("invalid base64 header")
			}
			payloadBytes, err := safeDecodeString(jwtParts[1])
			if err != nil {
				panic("invalid base64 payload")
			}

			header := JwtHeader{}
			err = json.Unmarshal(headerBytes, &header)
			if err != nil {
				fmt.Println("REJECTED bad_token")
				continue
			}
			header.Alg = strings.TrimSpace(header.Alg)
			// fmt.Println("header.Alg:", header.Alg)
			if !strings.HasPrefix(header.Alg, "HS") && !strings.HasPrefix(header.Alg, "RS") {
				fmt.Println("REJECTED bad_alg")
				continue
			}

			if strings.HasPrefix(header.Alg, "HS") {
				if len(hs256Keys) == 0 {
					fmt.Println("REJECTED alg_mismatch")
					continue
				}
				if header.Alg != "HS256" {
					fmt.Println("REJECTED bad_alg")
					continue
				}
				found := false
				for _, key := range hs256Keys {
					if len(key) == 0 {
						continue
					}
					expectedSig := signHs(headerBytes, payloadBytes, []byte(key))
					sigBytes, err := safeDecodeString(jwtParts[2])
					if err != nil {
						panic("invalid base64 signature")
					}
					if hmac.Equal(sigBytes, expectedSig) {
						fmt.Println("OK")
						found = true
						break
					}
				}
				if !found {
					fmt.Println("REJECTED bad_signature")
					continue
				}
			}

			if strings.HasPrefix(header.Alg, "RS") {
				if len(rsaPublicKeys) == 0 {
					fmt.Println("REJECTED alg_mismatch")
					continue
				}
				if header.Alg != "HS256" {
					fmt.Println("REJECTED bad_alg")
					continue
				}
			}

			// sig := jwtParts[2]
			// theSig := sign(headerBytes, payloadBytes, key)
			// sigBytes, err := safeDecodeString(jwtParts[2])
			// if err != nil {
			// 	panic("invalid base64 signature")
			// }
			// fmt.Println(expectedSig, sigBytes)

			// if hmac.Equal(sigBytes, expectedSig) {
			// 	fmt.Println("OK")
			// } else {
			// 	fmt.Println("BAD")
			// }
		}

	}
}
