// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
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

// ListenAndServe listens on the TCP network address addr and serves the handler.
// This function implements graceful shutdown when the passed ctx is done.
func ListenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: ReadHeaderTimeout,
	}

	egrp, ctx := errgroup.WithContext(ctx)

	egrp.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}

		return nil
	})

	egrp.Go(func() error {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %w", err)
		}

		return nil
	})

	if err := egrp.Wait(); err != nil {
		return fmt.Errorf("errgroup wait: %w", err)
	}

	return nil
}
