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
		ID: "Balance",
	}, nil
}

func (tenant *Tenant) Accept(ctx context.Context, t ot.Transferable) error {
	if t.GetOriginTenant() == t.GetDestinationTenant() {
		return nil
	}

	// money goes out from this tenant
	if t.GetOriginTenant() == tenant.ID() {
		if tenant.accounts[t.GetFromRef()] < t.GetValue() {
			return errors.New("not enough balance on the origin")
		}

		tenant.accounts[t.GetFromRef()] -= t.GetValue()

		return nil
	}

	// received balance
	tenant.accounts[t.GetToRef()] += t.GetValue()

	return nil
}

func (tenant *Tenant) AcceptRequest(ctx context.Context, t ot.Transferable) error {
	return nil
}

func (tenant *Tenant) Revert(ctx context.Context, t ot.Transferable) error {
	return nil
}
