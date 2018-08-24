package origin

import (
	"context"
	"errors"
	ot "github.com/petomalina/opentransaction"
)

type Tenant struct {
	accounts map[string]int
}

func (tenant *Tenant) Name() ot.TenantID {
	return "Origin"
}

func (tenant *Tenant) Accept(ctx context.Context, t ot.Transaction) error {
	if t.DestinationTenant() == tenant.Name() {
		return errors.New("origin can't accept any transactions, please see the documentation")
	}

	return nil
}

func (tenant *Tenant) AcceptRequest(ctx context.Context, t ot.Transaction) error {
	if t.DestinationTenant() == tenant.Name() {
		return errors.New("origin can't accept any transactions, please see the documentation")
	}

	return nil
}

func (tenant *Tenant) Revert(ctx context.Context, t ot.Transaction) error {
	return nil
}
