// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	tkhttp "github.com/hhromic/go-toolkit/http"
)

// HTTP headers used by the server package.
const (
	HeaderAuthorization      = "Authorization"
	HeaderXForwardedClientID = "X-Forwarded-Client-Id"
	HeaderXForwardedScope    = "X-Forwarded-Scope"
	HeaderXForwardedSubject  = "X-Forwarded-Subject"
)

const (
	// ShutdownTimeout is the maximum time to wait for the HTTP server to shutdown.
	ShutdownTimeout time.Duration = 30 * time.Second
	// ReadHeaderTimeout is the maximum time to wait for reading an HTTP request header.
	ReadHeaderTimeout time.Duration = 60 * time.Second
)

// Run listens on the TCP network address addr and serves the handler.
// This function implements graceful shutdown when the passed ctx is done.
func Run(ctx context.Context, addr string, handler http.Handler) error {
	srv := &http.Server{ //nolint:exhaustruct
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	if err := tkhttp.RunServer(ctx, srv, ShutdownTimeout); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

// Error replies to the request with the specified error message and HTTP code.
// It also logs the request remote address, error and code as a warning. The only
// exception are [http.StatusUnauthorized] codes, which are logged at the debug level.
func Error(writer http.ResponseWriter, request *http.Request, err string, code int) {
	http.Error(writer, err, code)

	logFn := slog.Warn
	if code == http.StatusUnauthorized {
		logFn = slog.Debug
	}

	logFn("request error", "addr", request.RemoteAddr, "err", err, "code", code)
}
