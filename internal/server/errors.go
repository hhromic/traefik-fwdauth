// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import "errors"

// Errors used by the server package.
var (
	// ErrRequestHeaderMissing is returned when a client request is missing a header.
	ErrMissingRequestHeader = errors.New("missing request header")

	// ErrUnsupportedAuthScheme is returned when a client request uses an unsupported authorization scheme.
	ErrUnsupportedAuthScheme = errors.New("unsupported authorization scheme")
	// ErrInactiveToken is returned when a client request presented an inactive token.
	ErrInactiveToken = errors.New("inactive token")
	// ErrInvalidClientID is returned when a client request presented an invalid client ID.
	ErrInvalidClientID = errors.New("invalid client ID")
)
