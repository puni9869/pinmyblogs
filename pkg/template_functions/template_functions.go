// Package template_functions provides custom template helper functions.
package template_functions //nolint:revive,stylecheck // existing package name

import (
	"fmt"
	"net/url"
	"time"
)

// DomainName extracts the hostname from a URL string.
func DomainName(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

// AvatarInitials returns the first character of the given text.
func AvatarInitials(text string) string {
	return fmt.Sprintf("%c", text[0])
}

// Add returns the sum of a and b.
func Add(a, b int) int {
	return a + b
}

// Sub returns the difference of a and b.
func Sub(a, b int) int {
	return a - b
}

// Asset returns a function that builds a versioned static asset path.
func Asset(version string) func(file string) string {
	return func(file string) string {
		return "/statics/" + file + "?v=" + version
	}
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
