package voucher

import (
	"context"
	ot "github.com/petomalina/opentransaction"
)

type Tenant struct{}

func (tenant *Tenant) Name() string {
	return "Voucher"
}

func (tenant *Tenant) Accept(ctx context.Context, t ot.Transaction) error {
	return nil
}

func (tenant *Tenant) AcceptRequest(ctx context.Context, t ot.Transaction) error {
	return nil
}

func (tenant *Tenant) Revert(ctx context.Context, t ot.Transaction) error {
	return nil
}
