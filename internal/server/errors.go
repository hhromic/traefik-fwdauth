// SPDX-FileCopyrightText: Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import "errors"

// Errors used by the server package.
var (
	// ErrRequestHeaderMissing is returned when a client request is missing a header.
	ErrMissingRequestHeader = errors.New("missing request header")
	// ErrUnsupportedAuthSyntax is returned when a client request uses an unsupported authorization syntax.
	ErrUnsupportedAuthSyntax = errors.New("unsupported authorization syntax")
)
