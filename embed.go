package pinmyblogs

import "embed"

//go:embed frontend/* templates/*
var Files embed.FS
