package origin

import (
	"context"
	"errors"
	ot "github.com/petomalina/opentransaction"
)

type Tenant struct{}

func (tenant *Tenant) ID(ctx context.Context, empty *ot.Empty) (*ot.TenantID, error) {
	return &ot.TenantID{
		ID: tenant.Name(),
	}, nil
}

func (tenant *Tenant) Name() string {
	return "[Tenant] Origin"
}

func (tenant *Tenant) Accept(ctx context.Context, t *ot.Transaction) (*ot.Metadata, error) {
	if t.GetDestinationTenant() == tenant.Name() {
		return nil, errors.New("origin can't accept any transactions, please see the documentation")
	}

	return &ot.Metadata{}, nil
}

func (tenant *Tenant) AcceptRequest(ctx context.Context, t *ot.Transaction) (*ot.Metadata, error) {
	if t.GetDestinationTenant() == tenant.Name() {
		return nil, errors.New("origin can't accept any transactions, please see the documentation")
	}

	return &ot.Metadata{}, nil
}

func (tenant *Tenant) Revert(ctx context.Context, t *ot.Transaction) (*ot.Metadata, error) {
	return &ot.Metadata{}, nil
}
