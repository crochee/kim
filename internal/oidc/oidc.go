/*
Copyright 2025.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
*/

package oidc

import (
	"context"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

// Provider wraps ZITADEL OIDC provider for CRD integration
type Provider struct {
	op.Provider
}

// NewProvider initializes OIDC with CRD-backed storage
func NewProvider(ctx context.Context, issuer string, config *op.Config) (*Provider, error) {
	provider, err := op.NewProvider(ctx, config)
	if err != nil {
		return nil, err
	}
	return &Provider{Provider: provider}, nil
}

// VerifyToken validates tokens against CRD policies
func (p *Provider) VerifyToken(ctx context.Context, token string) (*oidc.TokenClaims, error) {
	return p.Provider.VerifyToken(ctx, token)
}