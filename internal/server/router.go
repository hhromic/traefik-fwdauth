// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hhromic/traefik-fwdauth/v2/internal/client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// PatternAuthHandler is the path pattern to use for the auth handler.
	PatternAuthHandler = "/auth"
	// PatternMetricsHandler is the path pattern to use for the metrics handler.
	PatternMetricsHandler = "/metrics"
)

// NewRouter creates a top-level [http.Handler] router for the application.
func NewRouter(is *client.IntrospectionService) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Mount(PatternMetricsHandler, promhttp.Handler())
	r.Mount(PatternAuthHandler, AuthHandler(is))

	return r
}
