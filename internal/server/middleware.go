// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

//nolint:gochecknoglobals
var ctxKeyToken = &contextKey{"token"}

type contextKey struct {
	name string
}

// ExtractToken extracts client IDs from request query parameters.
func ExtractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		token, err := getToken(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)

			return
		}

		request = request.WithContext(context.WithValue(ctx, ctxKeyToken, token))

		next.ServeHTTP(writer, request)
	})
}

// TokenFromContext returns the token value stored in ctx, if any.
func TokenFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyToken).(string); ok {
		return v
	}

	return ""
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
