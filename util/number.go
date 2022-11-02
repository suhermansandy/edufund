package util

import (
	"crypto/sha256"
	"fmt"
)

func Encrypt(s string) string {
	encrypt := sha256.New()
	encrypt.Write([]byte(s))
	myString := fmt.Sprintf("%x", encrypt.Sum(nil))
	return myString
}

func IsFloat(n []byte) bool {
	if len(n) > 0 && n[0] == '-' {
		n = n[1:]
	}
	if len(n) == 0 {
		return false
	}
	var point bool
	for _, c := range n {
		if '0' <= c && c <= '9' {
			continue
		}
		if c == '.' && len(n) > 1 && !point {
			point = true
			continue
		}
		return false
	}
	return true
}
