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
	// AuthHandlerPattern is the path pattern to use for the auth handler.
	AuthHandlerPattern = "/auth"
	// MetricsHandlerPattern is the path pattern to use for the metrics handler.
	MetricsHandlerPattern = "/metrics"
)

// NewRouter creates a top-level http.Handler router for the application.
func NewRouter(is *client.IntrospectionService) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Mount(MetricsHandlerPattern, promhttp.Handler())
	r.Mount(AuthHandlerPattern, AuthHandler(is))

	return r
}
