// SPDX-FileCopyrightText: Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/twmb/go-cache/cache"
)

const (
	// FormFieldToken is the request form field used for providing a token.
	FormFieldToken = "token"
	// FormFieldTokenTypeHint is the request form field used for providing a token type hint.
	FormFieldTokenTypeHint = "token_type_hint"
)

// IntrospectionService is an OAuth 2.0 Token Introspection (RFC 7662) service for token validation.
type IntrospectionService struct {
	Client       *http.Client
	URL          url.URL
	ClientID     string
	ClientSecret string
	Cache        *cache.Cache[IntrospectionCacheKey, *IntrospectionResponse]
}

// IntrospectionCacheKey is the key used for caching introspection requests.
type IntrospectionCacheKey struct {
	Token         string
	TokenTypeHint string
}

// IntrospectionResponse is a response from the token introspection URL.
//
//nolint:tagliatelle
type IntrospectionResponse struct {
	Active   bool   `json:"active"`
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
	Subject  string `json:"sub"`
}

// NewIntrospectionCache creates a new cache to be used in an [IntrospectionService] instance.
func NewIntrospectionCache(
	ctx context.Context,
	expireAfter time.Duration,
) *cache.Cache[IntrospectionCacheKey, *IntrospectionResponse] {
	icache := cache.New[IntrospectionCacheKey, *IntrospectionResponse](
		cache.AutoCleanInterval(expireAfter/2), //nolint:mnd
		cache.MaxAge(expireAfter),
	)

	go func() {
		<-ctx.Done()
		icache.StopAutoClean()
	}()

	return icache
}

// Introspect performs token validation using token introspection.
func (s *IntrospectionService) Introspect(
	ctx context.Context,
	token, tokenTypeHint string,
) (*IntrospectionResponse, error) {
	cacheKey := IntrospectionCacheKey{
		Token:         token,
		TokenTypeHint: tokenTypeHint,
	}

	if ir, _, ks := s.Cache.TryGet(cacheKey); ks == cache.Hit {
		return ir, nil
	}

	introspectionURL := s.URL.String()

	form := &url.Values{}
	form.Set(FormFieldToken, cacheKey.Token)

	if cacheKey.TokenTypeHint != "" {
		form.Set(FormFieldTokenTypeHint, cacheKey.TokenTypeHint)
	}

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, introspectionURL, body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set(HeaderAccept, ContentTypeJSON)
	req.Header.Set(HeaderContentType, ContentTypeFormURLEncoded)
	req.SetBasicAuth(s.ClientID, s.ClientSecret)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client request: %w", err)
	}
	defer res.Body.Close() //nolint:errcheck

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %q", ErrBadResponse, res.Status)
	}

	var ires IntrospectionResponse
	if err := json.NewDecoder(res.Body).Decode(&ires); err != nil {
		return nil, fmt.Errorf("JSON decoder: %w", err)
	}

	s.Cache.Set(cacheKey, &ires)

	return &ires, nil
}
