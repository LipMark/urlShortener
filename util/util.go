package util

import (
	"log"
	"log/slog"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func SlogErr(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
