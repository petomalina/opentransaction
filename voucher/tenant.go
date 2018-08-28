package voucher

import (
	"context"
	ot "github.com/petomalina/opentransaction"
)

type Tenant struct{}

func (tenant *Tenant) ID(ctx context.Context, empty *ot.Empty) (*ot.TenantID, error) {
	return &ot.TenantID{
		ID: tenant.Name(),
	}, nil
}

func (tenant *Tenant) Name() string {
	return "[Tenant] Voucher"
}

func (tenant *Tenant) Accept(ctx context.Context, t *ot.Transaction) (*ot.Metadata, error) {
	return nil, nil
}

func (tenant *Tenant) AcceptRequest(ctx context.Context, t *ot.Transaction) (*ot.Metadata, error) {
	return nil, nil
}

func (tenant *Tenant) Revert(ctx context.Context, t *ot.Transaction) (*ot.Metadata, error) {
	return nil, nil
}
