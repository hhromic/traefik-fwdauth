// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/hhromic/traefik-fwdauth/v2/internal/client"
	"github.com/hhromic/traefik-fwdauth/v2/internal/metrics"
)

const (
	// QueryParamClientID is the request query parameter used for providing allowed client IDs.
	QueryParamClientID = "client_id"
	// QueryParamTokenTypeHint is the request query parameter used for providing a token type hint.
	QueryParamTokenTypeHint = "token_type_hint"
)

// AuthHandler is an [http.Handler] for authentication requests.
func AuthHandler(isrv *client.IntrospectionService) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		s := time.Now()
		defer func() {
			d := time.Since(s)
			slog.Debug("auth request completed", "duration", d)
			metrics.AuthRequestDuration.Observe(d.Seconds())
		}()

		token, err := getToken(request)
		if err != nil {
			handleUnauthorized(writer, err)

			return
		}

		tth := getTokenTypeHint(request)

		ires, err := isrv.Introspect(request.Context(), token, tth)
		if err != nil {
			handleErr(writer, fmt.Errorf("introspect: %w", err), http.StatusBadGateway)

			return
		}

		if !ires.Active {
			handleUnauthorized(writer, ErrInactiveToken)

			return
		}

		if !checkClientID(request, ires.ClientID) {
			handleUnauthorized(writer, ErrInvalidClientID)

			return
		}

		if ires.ClientID != "" {
			writer.Header().Set(HeaderXForwardedClientID, ires.ClientID)
		}

		if ires.Scope != "" {
			writer.Header().Set(HeaderXForwardedScope, ires.Scope)
		}

		if ires.Subject != "" {
			writer.Header().Set(HeaderXForwardedSubject, ires.Subject)
		}
	})
}

func getToken(r *http.Request) (string, error) {
	ahdr := r.Header.Get(HeaderAuthorization)
	if ahdr == "" {
		return "", fmt.Errorf("%w: %q", ErrMissingRequestHeader, HeaderAuthorization)
	}

	if len(ahdr) <= 7 || strings.ToUpper(ahdr[0:6]) != "BEARER" {
		return "", ErrUnsupportedAuthScheme
	}

	return ahdr[7:], nil
}

func getTokenTypeHint(r *http.Request) string {
	return r.URL.Query().Get(QueryParamTokenTypeHint)
}

func checkClientID(r *http.Request, cid string) bool {
	cids := r.URL.Query()[QueryParamClientID]
	if len(cids) == 0 {
		return true
	}

	for _, val := range cids {
		if val != "" && cid == val {
			return true
		}
	}

	return false
}

func handleUnauthorized(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
	slog.Debug("unauthorized auth request", "err", err)
	metrics.AuthRequestUnauthorized.Inc()
}

func handleErr(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
	slog.Error("auth handler error", "err", err, "status", status)
	metrics.AuthRequestErrors.Inc()
}
