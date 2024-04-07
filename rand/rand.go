package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("didn't read enough random bytes")
	}
	return b, nil
}

// string returns a random string using crypto/rand
// n is the number of bytes used to generate the random string
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("strings: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
