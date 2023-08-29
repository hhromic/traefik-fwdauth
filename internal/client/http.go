// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"net/http"
	"time"
)

// HTTP headers used by the client package.
const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"
)

// Content types used by the client package.
const (
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	ContentTypeJSON           = "application/json"
)

const (
	// ResponseHeaderTimeout is the maximum time to wait for reading an HTTP response header.
	ResponseHeaderTimeout = 60 * time.Second
)

// NewClient creates a new [http.Client] that uses a clone of [http.DefaultTransport]
// configured with a response header timeout.
func NewClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone() //nolint:forcetypeassert
	t.ResponseHeaderTimeout = ResponseHeaderTimeout

	c := &http.Client{ //nolint:exhaustruct,exhaustivestruct
		Transport: t,
	}

	return c
}
