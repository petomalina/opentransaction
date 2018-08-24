package balance

import (
	"context"
	"errors"
	ot "github.com/petomalina/opentransaction"
)

type Tenant struct {
	accounts map[string]int
}

func (tenant *Tenant) Name() ot.TenantID {
	return "Balance"
}

func (tenant *Tenant) Accept(ctx context.Context, t ot.Transaction) error {
	if t.OriginTenant() == t.DestinationTenant() {
		return nil
	}

	// money goes out from this tenant
	if t.OriginTenant() == tenant.Name() {
		if tenant.accounts[t.FromRef()] < t.Value() {
			return errors.New("not enough balance on the origin")
		}

		tenant.accounts[t.FromRef()] -= t.Value()

		return nil
	}

	// received balance
	tenant.accounts[t.ToRef()] += t.Value()

	return nil
}

func (tenant *Tenant) AcceptRequest(ctx context.Context, t ot.Transaction) error {
	return nil
}

func (tenant *Tenant) Revert(ctx context.Context, t ot.Transaction) error {
	return nil
}
