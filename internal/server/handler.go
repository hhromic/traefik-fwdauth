// Copyright 2023 Hugo Hromic
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net/http"

	"github.com/hhromic/traefik-fwdauth/v2/internal/client"
)

const (
	// QueryParamClientID is the request query parameter used for providing allowed client IDs.
	QueryParamClientID = "client_id"
	// QueryParamTokenTypeHint is the request query parameter used for providing a token type hint.
	QueryParamTokenTypeHint = "token_type_hint"
)

// AuthHandler is an [http.Handler] for authentication requests.
func AuthHandler(isrv *client.IntrospectionService) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		token := TokenFromContext(ctx)
		tth := request.URL.Query().Get(QueryParamTokenTypeHint)

		ires, err := isrv.Introspect(ctx, token, tth)
		if err != nil {
			Error(writer, request, "introspect: "+err.Error(), http.StatusBadGateway)

			return
		}

		if !ires.Active {
			Error(writer, request, "inactive token", http.StatusUnauthorized)

			return
		}

		if !isValidClientID(request, ires.ClientID) {
			Error(writer, request, "invalid client ID", http.StatusForbidden)

			return
		}

		if ires.ClientID != "" {
			writer.Header().Set(HeaderXForwardedClientID, ires.ClientID)
		}

		if ires.Scope != "" {
			writer.Header().Set(HeaderXForwardedScope, ires.Scope)
		}

		if ires.Subject != "" {
			writer.Header().Set(HeaderXForwardedSubject, ires.Subject)
		}
	})
}

func isValidClientID(r *http.Request, cid string) bool {
	cids := r.URL.Query()[QueryParamClientID]
	if len(cids) == 0 {
		return true
	}

	for _, val := range cids {
		if val != "" && cid == val {
			return true
		}
	}

	return false
}
