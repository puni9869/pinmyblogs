package utils

import (
	"crypto/sha256"
	"fmt"
	"net/mail"
)

// HashifySHA256 will convert any text into SHA256 hash
func HashifySHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// ValidateEmail will validate the email is in correct format or not
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
