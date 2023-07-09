// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/hhromic/traefik-fwdauth/v2/internal/buildinfo"
	"github.com/hhromic/traefik-fwdauth/v2/internal/client"
	"github.com/hhromic/traefik-fwdauth/v2/internal/logger"
	"github.com/hhromic/traefik-fwdauth/v2/internal/server"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/slog"

	_ "github.com/hhromic/traefik-fwdauth/v2/internal/metrics" // initialize collectors
)

type args struct {
	ListenAddress         string         `arg:"--listen-address,env:LISTEN_ADDRESS" default:":4181" placeholder:"ADDRESS" help:"listen address for the HTTP server"`
	OIDCIssuerURL         *url.URL       `arg:"--oidc-issuer-url,env:OIDC_ISSUER_URL" placeholder:"URL" help:"issuer URL for OIDC discovery"`
	IntrospectionEndpoint *url.URL       `arg:"--introspection-endpoint,env:INTROSPECTION_ENDPOINT" placeholder:"URL" help:"token introspection endpoint"`
	ClientID              string         `arg:"--client-id,required,env:CLIENT_ID" placeholder:"CLIENT_ID" help:"client ID for the token introspection endpoint"`
	ClientSecret          string         `arg:"--client-secret,env:CLIENT_SECRET" placeholder:"CLIENT_SECRET" help:"client secret for the token introspection endpoint"`
	ClientSecretFile      string         `arg:"--client-secret-file,env:CLIENT_SECRET_FILE" placeholder:"FILE" help:"file containing the client secret"`
	ExpireAfter           time.Duration  `arg:"--expire-after,env:EXPIRE_AFTER" default:"5m" placeholder:"DURATION" help:"time for expiring cached client requests"`
	LogHandler            logger.Handler `arg:"--log-handler,env:LOG_HANDLER" default:"text" placeholder:"HANDLER" help:"application logging handler"`
	LogLevel              slog.Level     `arg:"--log-level,env:LOG_LEVEL" default:"info" placeholder:"LEVEL" help:"application logging level"`
}

func main() {
	var args args
	p := arg.MustParse(&args)

	if args.OIDCIssuerURL == nil && args.IntrospectionEndpoint == nil {
		p.Fail("either --oidc-issuer-url or --introspection-endpoint is required")
	}

	if args.ClientSecret == "" && args.ClientSecretFile == "" {
		p.Fail("either --client-secret or --client-secret-file is required")
	}

	if err := logger.SlogSetDefault(os.Stderr, args.LogHandler, args.LogLevel); err != nil {
		panic(err)
	}

	if _, err := maxprocs.Set(); err != nil {
		slog.Warn("failed to set GOMAXPROCS", "err", err)
	}

	slog.Info("starting",
		"version", buildinfo.Version,
		"goversion", buildinfo.GoVersion,
		"gitcommit", buildinfo.GitCommit,
		"gitbranch", buildinfo.GitBranch,
		"builddate", buildinfo.BuildDate,
		"gomaxprocs", runtime.GOMAXPROCS(0),
	)

	if args.ClientSecretFile != "" {
		data, err := os.ReadFile(args.ClientSecretFile)
		if err != nil {
			slog.Error("error reading client secret from file", "err", err)
			os.Exit(1)
		}

		args.ClientSecret = strings.TrimRight(string(data), "\r\n")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	clnt := client.NewClient()

	if args.OIDCIssuerURL != nil {
		ds := &client.OIDCDiscoveryService{
			Client:    clnt,
			IssuerURL: *args.OIDCIssuerURL,
		}

		endpoint, err := ds.DiscoverIntrospection(ctx)
		if err != nil {
			slog.Error("OIDC discovery failed", "err", err)
			os.Exit(1)
		}

		slog.Info("OIDC discovery completed", "introspection_endpoint", endpoint)
		args.IntrospectionEndpoint = endpoint
	}

	isrv := &client.IntrospectionService{
		Client:       clnt,
		URL:          *args.IntrospectionEndpoint,
		ClientID:     args.ClientID,
		ClientSecret: args.ClientSecret,
		Cache:        client.NewIntrospectionCache(ctx, args.ExpireAfter),
	}

	r := server.NewRouter(isrv)

	slog.Info("starting HTTP server", "addr", args.ListenAddress)

	if err := server.ListenAndServe(ctx, args.ListenAddress, r); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("error running HTTP server", "err", err)
	} else {
		slog.Info("finished")
	}
}
