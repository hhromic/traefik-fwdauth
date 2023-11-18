// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net/http"

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

// NewServeMux creates a top-level request multiplexer for the application.
func NewServeMux(isrv *client.IntrospectionService) *http.ServeMux {
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

	m := http.NewServeMux()
	m.Handle(PatternMetricsHandler, promhttp.Handler())
	m.Handle(PatternAuthHandler, ahandler)

	return m
}
