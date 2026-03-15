package api

import (
	"log/slog"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func logWithRef(err error, op string) string {
	refCode, _ := gonanoid.New(21)
	slog.Error("internal server error",
		"operation", op,
		"ref", refCode,
		"error", err,
	)

	return refCode
}
