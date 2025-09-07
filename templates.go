package assets

import "embed"

//go:embed templates/*
var Templates embed.FS

// - embed.FS = virtual filesystem in binary
// - go:embed = selects which files to include
