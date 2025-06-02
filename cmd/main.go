// SPDX-FileCopyrightText: Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	tkslog "github.com/hhromic/go-toolkit/slog"
	"github.com/hhromic/traefik-fwdauth/v2/internal/buildinfo"
	"github.com/hhromic/traefik-fwdauth/v2/internal/client"
	_ "github.com/hhromic/traefik-fwdauth/v2/internal/metrics" // initialize collectors
	"github.com/hhromic/traefik-fwdauth/v2/internal/server"
	"go.uber.org/automaxprocs/maxprocs"
)

//nolint:lll,tagalign
type args struct {
	ListenAddress         string         `arg:"--listen-address,env:LISTEN_ADDRESS" default:":4181" placeholder:"ADDRESS" help:"listen address for the HTTP server"`
	OIDCIssuerURL         *url.URL       `arg:"--oidc-issuer-url,env:OIDC_ISSUER_URL" placeholder:"URL" help:"issuer URL for OIDC discovery"`
	IntrospectionEndpoint *url.URL       `arg:"--introspection-endpoint,env:INTROSPECTION_ENDPOINT" placeholder:"URL" help:"token introspection endpoint"`
	ClientID              string         `arg:"--client-id,required,env:CLIENT_ID" placeholder:"CLIENT_ID" help:"client ID for the token introspection endpoint"`
	ClientSecret          string         `arg:"--client-secret,env:CLIENT_SECRET" placeholder:"CLIENT_SECRET" help:"client secret for the token introspection endpoint"`
	ClientSecretFile      string         `arg:"--client-secret-file,env:CLIENT_SECRET_FILE" placeholder:"FILE" help:"file containing the client secret"`
	ExpireAfter           time.Duration  `arg:"--expire-after,env:EXPIRE_AFTER" default:"5m" placeholder:"DURATION" help:"time for expiring cached client requests"`
	LogHandler            tkslog.Handler `arg:"--log-handler,env:LOG_HANDLER" default:"auto" placeholder:"HANDLER" help:"application logging handler"`
	LogLevel              slog.Level     `arg:"--log-level,env:LOG_LEVEL" default:"info" placeholder:"LEVEL" help:"application logging level"`
}

func (args) Description() string {
	return "Traefik forward auth service version " + buildinfo.Version +
		" (git:" + buildinfo.GitBranch + "/" + buildinfo.GitCommit + ")"
}

func main() {
	var args args
	parser := arg.MustParse(&args)

	if args.OIDCIssuerURL == nil && args.IntrospectionEndpoint == nil {
		parser.Fail("either --oidc-issuer-url or --introspection-endpoint is required")
	}

	if args.ClientSecret == "" && args.ClientSecretFile == "" {
		parser.Fail("either --client-secret or --client-secret-file is required")
	}

	slog.SetDefault(tkslog.NewSlogLogger(os.Stderr, args.LogHandler, args.LogLevel))

	if err := appMain(args); err != nil {
		slog.Error("application error", "err", err)
		os.Exit(1)
	}
}

func appMain(args args) error { //nolint:funlen
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
			return fmt.Errorf("error reading client secret from file: %w", err)
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
			return fmt.Errorf("OIDC discovery failed: %w", err)
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

	m := server.NewServeMux(isrv)

	slog.Info("starting HTTP server", "addr", args.ListenAddress)

	err := server.Run(ctx, args.ListenAddress, m)
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("run: %w", err)
	}

	slog.Info("finished")

	return nil
}
