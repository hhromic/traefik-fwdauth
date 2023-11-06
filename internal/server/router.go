// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hhromic/traefik-fwdauth/v2/internal/client"
	"github.com/hhromic/traefik-fwdauth/v2/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// PatternAuthHandler is the path pattern to use for the auth handler.
	PatternAuthHandler = "/auth"
	// PatternMetricsHandler is the path pattern to use for the metrics handler.
	PatternMetricsHandler = "/metrics"
)

// NewRouter creates a top-level [http.Handler] router for the application.
func NewRouter(isrv *client.IntrospectionService) http.Handler {
	ahandler := promhttp.InstrumentHandlerInFlight(
		metrics.AuthInFlightRequests,
		promhttp.InstrumentHandlerDuration(
			metrics.AuthRequestDuration,
			promhttp.InstrumentHandlerCounter(
				metrics.AuthRequestsTotal,
				ExtractToken(AuthHandler(isrv)),
			),
		),
	)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Mount(PatternMetricsHandler, promhttp.Handler())
	r.Mount(PatternAuthHandler, ahandler)

	return r
}
