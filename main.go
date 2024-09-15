package main

import (
	"fmt"
	internalcfg "urlShortener/internal/internalconfig"
)

func main() {
	cfg := internalcfg.MustLoad()
	// config
	fmt.Println(cfg)
	// logger

	// stor

	// router

	// server
}
