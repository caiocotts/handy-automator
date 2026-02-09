package main

import (
	"embed"
)

//go:embed banner.txt
var banner embed.FS

func Banner() string {
	data, _ := banner.ReadFile("banner.txt")
	return string(data)
}
