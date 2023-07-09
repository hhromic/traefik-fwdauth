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

// NewClient creates a new http.Client that uses http.DefaultTransport with a
// configured response header timeout.
func NewClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.ResponseHeaderTimeout = ResponseHeaderTimeout

	c := &http.Client{
		Transport: t,
	}

	return c
}
