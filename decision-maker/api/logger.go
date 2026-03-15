package api

import (
	"log"
	"log/slog"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func logWithRef(err error, op string) string {
	refCode, err := gonanoid.New(21)
	if err != nil {
		log.Fatal("error: failed to generate reference code for log message: ", err)
	}

	slog.Error("internal server error",
		"operation", op,
		"ref", refCode,
		"error", err,
	)

	return refCode
}
