package services

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func hashKey(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{'|'})
	}
	return hex.EncodeToString(h.Sum(nil))
}
