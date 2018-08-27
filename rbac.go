package opentransaction

import (
	"fmt"
	"github.com/pkg/errors"
)

type RBACDefaultPolicy string

const (
	RBACDefaultOpen   RBACDefaultPolicy = "[rbac] open"
	RBACDefaultClosed                   = "[rbac] closed"
)

var (
	MissingOriginRBACPolicyErr      = errors.New("missing origin RBAC policy for the transaction")
	MissingDestinationRBACPolicyErr = errors.New("missing destination RBAC policy for the transaction")
)

type RBAC struct {
	defaultPolicy RBACDefaultPolicy

	policyMap map[string][]string
}

func NewRBAC(opts ...RBACOption) *RBAC {
	rbac := &RBAC{
		defaultPolicy: RBACDefaultOpen,
		policyMap:     map[string][]string{},
	}

	for _, opt := range opts {
		opt(rbac)
	}

	return rbac
}

type RBACOption func(rbac *RBAC)

func WithClosedPolicy() RBACOption {
	return func(rbac *RBAC) {
		rbac.defaultPolicy = RBACDefaultClosed
	}
}

func WithPolicy(origin, destination string) RBACOption {
	return func(rbac *RBAC) {
		if _, ok := rbac.policyMap[origin]; !ok {
			rbac.policyMap[origin] = make([]string, 1)
		}

		rbac.policyMap[origin] = append(rbac.policyMap[origin], destination)
	}
}

func (rbac *RBAC) Enforce(t Transaction) error {
	fmt.Println(rbac.policyMap, t.GetOriginTenant())
	// open policy will automatically authorize any requests
	if rbac.defaultPolicy == RBACDefaultOpen {
		return nil
	}

	policies, ok := rbac.policyMap[t.GetOriginTenant()]
	if !ok {
		return errors.Wrap(MissingOriginRBACPolicyErr, "enforce failed for origin RBAC policy on: "+string(t.GetOriginTenant()))
	}

	for _, p := range policies {
		// policy found
		if p == t.GetDestinationTenant() {
			return nil
		}
	}

	return errors.Wrap(MissingDestinationRBACPolicyErr, "enforce failed for destination RBAC policy on: "+string(t.GetDestinationTenant()))
}
