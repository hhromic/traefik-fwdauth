// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

// Handler represents a supported slog handler.
type Handler int

// Supported slog handlers.
const (
	// HandlerText is an [slog.TextHandler] which outputs logs in key=value format.
	HandlerText Handler = iota
	// HandlerJSON is an [slog.JSONHandler] which outputs logs in standard JSON format.
	HandlerJSON
	// HandlerTint is an [slog.Handler] which outputs colorized logs in key=value format.
	HandlerTint
	// HandlerAuto uses [HandlerTint] if the output writer is a terminal or [HandlerText] otherwise.
	HandlerAuto
)

// Errors used by the logger package.
var (
	// ErrUnknownHandlerName is returned when an unknown slog handler name is used.
	ErrUnknownHandlerName = errors.New("unknown handler name")
)

// String returns a name for the slog handler.
func (h Handler) String() string {
	switch h {
	case HandlerText:
		return "text"
	case HandlerJSON:
		return "json"
	case HandlerTint:
		return "tint"
	case HandlerAuto:
		return "auto"
	default:
		return ""
	}
}

// MarshalText implements [encoding.TextMarshaler] by calling [Handler.String].
func (h Handler) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler].
// It accepts any string produced by [Handler.MarshalText], ignoring case.
func (h *Handler) UnmarshalText(b []byte) error {
	str := string(b)
	switch strings.ToLower(str) {
	case "text":
		*h = HandlerText
	case "json":
		*h = HandlerJSON
	case "tint":
		*h = HandlerTint
	case "auto":
		*h = HandlerAuto
	default:
		return fmt.Errorf("%w: %q", ErrUnknownHandlerName, str)
	}

	return nil
}

// SlogSetDefault sets the global slog logger to output to writer, using the specified log handler
// and the specified minimum logging level. This function also renames the built-in attribute
// [slog.TimeKey] to "ts" for shorter logs output.
func SlogSetDefault(writer io.Writer, handler Handler, level slog.Leveler) error {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Key = "ts"
			}

			return a
		},
	}

	if handler == HandlerAuto {
		handler = HandlerText

		if f, ok := writer.(*os.File); ok {
			if isatty.IsTerminal(f.Fd()) {
				handler = HandlerTint
			}
		}
	}

	switch handler {
	case HandlerText:
		slog.SetDefault(slog.New(slog.NewTextHandler(writer, opts)))
	case HandlerJSON:
		slog.SetDefault(slog.New(slog.NewJSONHandler(writer, opts)))
	case HandlerTint:
		topts := &tint.Options{ //nolint:exhaustruct,exhaustivestruct
			AddSource: opts.AddSource,
			Level:     opts.Level,
		}
		slog.SetDefault(slog.New(tint.NewHandler(writer, topts)))
	case HandlerAuto:
	}

	return nil
}
