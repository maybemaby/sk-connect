package api

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
)

type LoggingFormat string

const (
	JSONFormat LoggingFormat = "json"
	TextFormat LoggingFormat = "text"
)

func BootstrapLogger(level slog.Level, format LoggingFormat, colorize bool) *slog.Logger {
	handlers := []slog.Handler{}

	if format == JSONFormat {
		handlers = append(handlers, slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		}))
	} else if format == TextFormat {

		if colorize {
			handlers = append(handlers, tint.NewHandler(os.Stderr, &tint.Options{
				Level:      level,
				TimeFormat: time.Kitchen,
			}))
		} else {
			handlers = append(handlers, slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: level,
			}))
		}
	}

	handler := slogmulti.Fanout(handlers...)

	return slog.New(handler)
}
