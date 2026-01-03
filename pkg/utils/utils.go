package utils

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

// HashifySHA256 will convert any text into SHA256 hash
func HashifySHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// HashifyBCrypt will convert any text into SHA256 hash
func HashifyBCrypt(text string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(h)
}

func CompareBCrypt(hashed string, plainText string) error {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashed),
		[]byte(plainText),
	); err != nil {
		return fmt.Errorf("compare bcrypt hash: %w", err)
	}

	return nil
}

// ValidateEmail will validate the email is in correct format or not
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
