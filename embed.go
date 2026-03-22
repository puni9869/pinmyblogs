// Package pinmyblogs embeds static frontend assets and HTML templates.
package pinmyblogs

import "embed"

// Files holds the embedded frontend and template assets.
//
//go:embed frontend/* templates/*
var Files embed.FS
