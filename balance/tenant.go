package balance

import (
	"context"
	"errors"
	ot "github.com/petomalina/opentransaction"
)

type Tenant struct {
	accounts map[string]int64
}

func (tenant *Tenant) ID(ctx context.Context, empty *ot.Empty) (*ot.TenantID, error) {
	return &ot.TenantID{
		ID: tenant.Name(),
	}, nil
}

func (tenant *Tenant) Name() string {
	return "[Tenant] Balance"
}

func (tenant *Tenant) Accept(ctx context.Context, t ot.Transaction) (*ot.Metadata, error) {
	if t.GetOriginTenant() == t.GetDestinationTenant() {
		return &ot.Metadata{}, nil
	}

	// money goes out from this tenant
	if t.GetOriginTenant() == tenant.Name() {
		if tenant.accounts[t.GetFromRef()] < t.GetValue() {
			return nil, errors.New("not enough balance on the origin")
		}

		tenant.accounts[t.GetFromRef()] -= t.GetValue()

		return &ot.Metadata{}, nil
	}

	// received balance
	tenant.accounts[t.GetToRef()] += t.GetValue()

	return &ot.Metadata{}, nil
}

func (tenant *Tenant) AcceptRequest(ctx context.Context, t ot.Transaction) (*ot.Metadata, error) {
	return &ot.Metadata{}, nil
}

func (tenant *Tenant) Revert(ctx context.Context, t ot.Transaction) (*ot.Metadata, error) {
	return &ot.Metadata{}, nil
}
