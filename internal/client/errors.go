// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package client

import "errors"

// Errors used by the client package.
var (
	// ErrBadResponse is returned when a bad server response is received.
	ErrBadResponse = errors.New("bad server response")

	// ErrDiscoveryMetadataMissing is returned when OIDC discovery metadata is missing.
	ErrDiscoveryMetadataMissing = errors.New("missing OIDC discovery metadata")
)
