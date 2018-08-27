package opentransaction

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/gofrs/uuid"
)

type Transferable interface {
	GetID() string

	GetOriginTenant() string
	GetDestinationTenant() string

	GetFromRef() string
	GetToRef() string

	GetValue() int64
}

func NewTransaction(from, to, originT, destinationT string, value int64) *Transaction {
	id, _ := uuid.NewV4()

	return &Transaction{
		ID:                id.String(),
		OriginTenant:      originT,
		DestinationTenant: destinationT,
		FromRef:           from,
		ToRef:             to,
		Value:             value,
	}
}

func NewRandomTransaction() *Transaction {
	id, _ := uuid.NewV4()

	return &Transaction{
		ID:                id.String(),
		OriginTenant:      "Balance",
		DestinationTenant: "Bank",
		FromRef:           "A",
		ToRef:             "B",
		Value:             100,
	}
}

type Core struct {
	tenants map[string]TenantServer

	rbac *RBAC
}

func NewCore(rbacOpts ...RBACOption) *Core {
	return &Core{
		tenants: make(map[string]TenantServer),
		rbac:    NewRBAC(rbacOpts...),
	}
}

func (c *Core) RegisterTenant(t TenantServer) error {
	c.tenants[t.ID()] = t

	return nil
}

func (c *Core) Send(tt ...Transferable) error {
	// enforce RBAC policies
	for _, t := range tt {
		if err := c.rbac.Enforce(t); err != nil {
			return errors.Wrap(err, "an RBAC policy failed for the transaction"+fmt.Sprintf("%+v", t))
		}
	}

	// send transactions
	var failIndex = 0
	var t Transferable
	var err error

	for failIndex, t = range tt {
		fmt.Printf("Sending transaction %+v\n", t)

		if err = c.tenants[t.GetOriginTenant()].Accept(context.Background(), t); err != nil {
			break
		}

		if err = c.tenants[t.GetDestinationTenant()].Accept(context.Background(), t); err != nil {
			c.tenants[t.GetOriginTenant()].Revert(context.Background(), t)
			break
		}
	}

	// we will need to revert all transactions until failIndex if it failed
	if err != nil {
		for i, t := range tt {
			// done on failIndex, we don't want to revert anything that was not processed
			if i == failIndex {
				break
			}

			// revert both origin and destination for each transaction
			c.tenants[t.GetOriginTenant()].Revert(context.Background(), t)
			c.tenants[t.GetDestinationTenant()].Revert(context.Background(), t)
		}

		return errors.New("an error occured during the transaction (reverted): " + err.Error())
	}

	return nil
}

func (c *Core) SendRequest(t Transferable) error {
	return nil
}
