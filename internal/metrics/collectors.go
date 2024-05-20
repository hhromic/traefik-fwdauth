// SPDX-FileCopyrightText: Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"github.com/hhromic/traefik-fwdauth/v2/internal/buildinfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// Namespace is the metrics namespace for the application.
	Namespace = "fwdauth"
)

// BuildInfo is the collector for build information of the application.
//
//nolint:gochecknoglobals
var BuildInfo = promauto.NewGaugeFunc(
	prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "build",
		Name:      "info",
		Help: "A metric with a constant '1' value labeled by version, goversion, gitcommit, " +
			"gitbranch, builddate from which the application was built.",
		ConstLabels: prometheus.Labels{
			"version":   buildinfo.Version,
			"goversion": buildinfo.GoVersion,
			"gitcommit": buildinfo.GitCommit,
			"gitbranch": buildinfo.GitBranch,
			"builddate": buildinfo.BuildDate,
		},
	},
	func() float64 { return 1 },
)

// AuthInFlightRequests is the collector for the number of auth requests currently being served.
//
//nolint:gochecknoglobals
var AuthInFlightRequests = promauto.NewGauge(
	prometheus.GaugeOpts{
		Namespace:   Namespace,
		Subsystem:   "auth",
		Name:        "in_flight_requests",
		Help:        "Number of auth requests currently being served in the Traefik Forward Auth service.",
		ConstLabels: prometheus.Labels{},
	},
)

// AuthRequestsTotal is the collector for the total number of auth requests.
//
//nolint:gochecknoglobals
var AuthRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace:   Namespace,
		Subsystem:   "auth",
		Name:        "requests_total",
		Help:        "Total number of auth requests in the Traefik Forward Auth service.",
		ConstLabels: prometheus.Labels{},
	},
	[]string{"code"},
)

// AuthRequestDuration is the collector for the distribution of auth request durations.
//
//nolint:exhaustruct,gochecknoglobals
var AuthRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace:   Namespace,
		Subsystem:   "auth",
		Name:        "request_duration_seconds",
		Help:        "Distribution of auth request durations in the Traefik Forward Auth service.",
		Buckets:     []float64{.1, .2, .4, 1, 3, 8, 20, 60, 120},
		ConstLabels: prometheus.Labels{},
	},
	[]string{"code"},
)
