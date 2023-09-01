// Copyright 2023 Hugo Hromic
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

// AuthRequestDuration is the collector for the distribution of auth request durations.
//
//nolint:exhaustruct,gochecknoglobals
var AuthRequestDuration = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Namespace:   Namespace,
		Subsystem:   "auth",
		Name:        "request_duration_seconds",
		Help:        "Distribution of auth request durations in the Traefik Forward Auth service.",
		ConstLabels: prometheus.Labels{},
	},
)

// AuthRequestErrors is the collector for the total number of auth request errors.
//
//nolint:gochecknoglobals
var AuthRequestErrors = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace:   Namespace,
		Subsystem:   "auth",
		Name:        "request_errors_total",
		Help:        "Total number of auth request errors in the Traefik Forward Auth service.",
		ConstLabels: prometheus.Labels{},
	},
)

// AuthRequestUnauthorized is the collector for the total number of unauthorized auth requests.
//
//nolint:gochecknoglobals
var AuthRequestUnauthorized = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace:   Namespace,
		Subsystem:   "auth",
		Name:        "request_unauthorized_total",
		Help:        "Total number of unauthorized auth requests in the Traefik Forward Auth service.",
		ConstLabels: prometheus.Labels{},
	},
)
