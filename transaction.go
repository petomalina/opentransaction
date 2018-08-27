package opentransaction

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/gofrs/uuid"
)

type TenantID string

type Transaction interface {
	ID() string

	OriginTenant() TenantID
	DestinationTenant() TenantID

	FromRef() string
	ToRef() string

	Value() int
}

type SimpleTransaction struct {
	id, fromRef, toRef              string
	originTenant, destinationTenant TenantID
	value                           int
}

func NewSimpleTransaction(from, to string, originT, destinationT TenantID, value int) *SimpleTransaction {
	id, _ := uuid.NewV4()

	return &SimpleTransaction{
		id:                id,
		originTenant:      originT,
		destinationTenant: destinationT,
		fromRef:           from,
		toRef:             to,
		value:             value,
	}
}

func (t *SimpleTransaction) ID() string {
	return t.id
}

func (t *SimpleTransaction) OriginTenant() TenantID {
	return t.originTenant
}

func (t *SimpleTransaction) DestinationTenant() TenantID {
	return t.destinationTenant
}

func (t *SimpleTransaction) FromRef() string {
	return t.fromRef
}

func (t *SimpleTransaction) ToRef() string {
	return t.toRef
}

func (t *SimpleTransaction) Value() int {
	return t.value
}

type Tenant interface {
	ID() TenantID

	Accept(ctx context.Context, t Transaction) error
	AcceptRequest(ctx context.Context, t Transaction) error

	Revert(ctx context.Context, t Transaction) error
}

type Core struct {
	tenants map[TenantID]Tenant

	rbac *RBAC
}

func NewCore(rbacOpts ...RBACOption) *Core {
	return &Core{
		tenants: make(map[TenantID]Tenant),
		rbac:    NewRBAC(rbacOpts...),
	}
}

func (c *Core) RegisterTenant(t Tenant) error {
	c.tenants[t.ID()] = t

	return nil
}

func (c *Core) Send(tt ...Transaction) error {
	// enforce RBAC policies
	for _, t := range tt {
		if err := c.rbac.Enforce(t); err != nil {
			return errors.Wrap(err, "an RBAC policy failed for the transaction"+fmt.Sprintf("%+v", t))
		}
	}

	// send transactions
	var failIndex = 0
	var t Transaction
	var err error

	for failIndex, t = range tt {
		fmt.Printf("Sending transaction %+v\n", t)

		if err = c.tenants[t.OriginTenant()].Accept(context.Background(), t); err != nil {
			break
		}

		if err = c.tenants[t.DestinationTenant()].Accept(context.Background(), t); err != nil {
			c.tenants[t.OriginTenant()].Revert(context.Background(), t)
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
			c.tenants[t.OriginTenant()].Revert(context.Background(), t)
			c.tenants[t.DestinationTenant()].Revert(context.Background(), t)
		}

		return errors.New("an error occured during the transaction (reverted): " + err.Error())
	}

	return nil
}

func (c *Core) SendRequest(t Transaction) error {
	return nil
}
