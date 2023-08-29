// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// PathWellKnownOpenIDConfiguration is the path in the issuer for performing OIDC discovery.
	PathWellKnownOpenIDConfiguration = "/.well-known/openid-configuration"
)

// OIDCDiscoveryService is an OIDC discovery service for obtaining OIDC resource metadata.
type OIDCDiscoveryService struct {
	Client    *http.Client
	IssuerURL url.URL
}

// OIDCDiscoveryResponse is a response from the OIDC discovery URL.
//
//nolint:tagliatelle
type OIDCDiscoveryResponse struct {
	// IntrospectionEndpoint is the URL for OAuth 2.0 Token Introspection (RFC 7662).
	IntrospectionEndpoint string `json:"introspection_endpoint"`
}

// Discover computes the OIDC discovery URL from the configured issuer URL.
func (s *OIDCDiscoveryService) Discover(ctx context.Context) (*OIDCDiscoveryResponse, error) {
	discoveryURL := s.IssuerURL.JoinPath(PathWellKnownOpenIDConfiguration).String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set(HeaderAccept, ContentTypeJSON)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %q", ErrBadResponse, res.Status)
	}

	var odr OIDCDiscoveryResponse
	if err := json.NewDecoder(res.Body).Decode(&odr); err != nil {
		return nil, fmt.Errorf("JSON decoder: %w", err)
	}

	return &odr, nil
}

// OIDCDiscoverIntrospection discovers an introspection URL using OIDC discovery.
func (s *OIDCDiscoveryService) DiscoverIntrospection(ctx context.Context) (*url.URL, error) {
	odr, err := s.Discover(ctx)
	if err != nil {
		return nil, fmt.Errorf("discover: %w", err)
	}

	if odr.IntrospectionEndpoint == "" {
		return nil, fmt.Errorf("%w: introspection_endpoint", ErrDiscoveryMetadataMissing)
	}

	u, err := url.ParseRequestURI(odr.IntrospectionEndpoint)
	if err != nil {
		return nil, fmt.Errorf("URL parser: %w", err)
	}

	return u, nil
}
