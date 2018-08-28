package opentransaction

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/gofrs/uuid"
)

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
	tenants map[string]TenantClient

	rbac *RBAC
}

func NewCore(rbacOpts ...RBACOption) *Core {
	return &Core{
		tenants: make(map[string]TenantClient),
		rbac:    NewRBAC(rbacOpts...),
	}
}

func (c *Core) RegisterTenant(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	tenantClient := NewTenantClient(conn)
	id, err := tenantClient.ID(context.Background(), nil)
	if err != nil {
		return err
	}

	c.tenants[id.ID] = tenantClient

	return nil
}

func (c *Core) Shutdown() error {
	return nil
}

func (c *Core) Send(ctx context.Context, tt *Transactions) (*Metadata, error) {
	// enforce RBAC policies
	for _, t := range tt.Transactions {
		if err := c.rbac.Enforce(t); err != nil {
			return nil, errors.Wrap(err, "an RBAC policy failed for the transaction"+fmt.Sprintf("%+v", t))
		}
	}

	// send transactions
	var failIndex = 0
	var t *Transaction
	var err error

	for failIndex, t = range tt.Transactions {
		fmt.Printf("Sending transaction %+v\n", t)

		if _, err = c.tenants[t.GetOriginTenant()].Accept(context.Background(), t); err != nil {
			break
		}

		if _, err = c.tenants[t.GetDestinationTenant()].Accept(context.Background(), t); err != nil {
			c.tenants[t.GetOriginTenant()].Revert(context.Background(), t)
			break
		}
	}

	// we will need to revert all transactions until failIndex if it failed
	if err != nil {
		for i, t := range tt.Transactions {
			// done on failIndex, we don't want to revert anything that was not processed
			if i == failIndex {
				break
			}

			// revert both origin and destination for each transaction
			c.tenants[t.GetOriginTenant()].Revert(context.Background(), t)
			c.tenants[t.GetDestinationTenant()].Revert(context.Background(), t)
		}

		return nil, errors.New("an error occured during the transaction (reverted): " + err.Error())
	}

	return &Metadata{}, nil
}

func (c *Core) SendRequest(t Transaction) error {
	return nil
}
