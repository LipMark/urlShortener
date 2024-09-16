package util

import (
	"log"
	"log/slog"
)

// CheckError for unexpected errors Os.Exit(1)
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// SlogErr handling slog errors
func SlogErr(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
