package origin

import (
	"context"
	"errors"
	ot "github.com/petomalina/opentransaction"
)

type Tenant struct{}

func (tenant *Tenant) Name() string {
	return "Origin"
}

func (tenant *Tenant) Accept(ctx context.Context, t ot.Transferable) error {
	if t.GetDestinationTenant() == tenant.Name() {
		return errors.New("origin can't accept any transactions, please see the documentation")
	}

	return nil
}

func (tenant *Tenant) AcceptRequest(ctx context.Context, t ot.Transferable) error {
	if t.GetDestinationTenant() == tenant.Name() {
		return errors.New("origin can't accept any transactions, please see the documentation")
	}

	return nil
}

func (tenant *Tenant) Revert(ctx context.Context, t ot.Transferable) error {
	return nil
}
