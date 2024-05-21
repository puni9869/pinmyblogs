package utils

import (
	"crypto/sha256"
	"fmt"
)

// HashifySHA256 will convert any text into SHA256 hash
func HashifySHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
