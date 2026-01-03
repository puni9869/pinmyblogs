package utils

import (
	"crypto/sha256"
	"fmt"
	"net/mail"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func DomainName(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// FormatRelativeTime returns a string like "3 hours ago" or "2 days ago"
func FormatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "unknown time"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%d months ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d years ago", int(d.Hours()/(24*365)))
	}
}

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
