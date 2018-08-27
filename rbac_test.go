package opentransaction

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RBACSuite struct {
	suite.Suite
}

func (s *RBACSuite) TestEnforce() {
	candidates := []struct {
		options      []RBACOption
		transactions map[Transferable]error
	}{
		{
			options: []RBACOption{},
			transactions: map[Transferable]error{
				NewTransaction("A", "B", "Balance", "Balance", 100): nil,
				NewTransaction("B", "A", "Balance", "Balance", 100): nil,
			},
		},
		{
			options: []RBACOption{
				WithClosedPolicy(),
			},
			transactions: map[Transferable]error{
				NewTransaction("A", "B", "Balance", "Balance", 100): MissingOriginRBACPolicyErr,
				NewTransaction("B", "A", "Balance", "Balance", 100): MissingOriginRBACPolicyErr,
			},
		},
		{
			options: []RBACOption{
				WithClosedPolicy(),
				WithPolicy("Balance", "Balance"),
			},
			transactions: map[Transferable]error{
				NewTransaction("A", "B", "Balance", "Balance", 100): nil,
				NewTransaction("B", "A", "Balance", "Balance", 100): nil,
				NewTransaction("B", "A", "Balance", "Bank", 100):    MissingDestinationRBACPolicyErr,
				NewTransaction("B", "A", "Bank", "Balance", 100):    MissingOriginRBACPolicyErr,
			},
		},
	}

	for _, c := range candidates {
		rbac := NewRBAC(c.options...)

		for t, errExpect := range c.transactions {
			err := rbac.Enforce(t)
			s.Equal(errExpect, errors.Cause(err), fmt.Sprintf("%+v", c))
		}
	}
}

func TestRBACSuite(t *testing.T) {
	suite.Run(t, &RBACSuite{})
}
