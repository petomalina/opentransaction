package origin

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TenantSuite struct {
	suite.Suite
	tenant *Tenant
}

func (s *TenantSuite) SetupTest() {
	s.tenant = &Tenant{}
}

func (s *TenantSuite) TestAccept() {
	s.tenant.Accept(context.Background(), nil)
}

func (s *TenantSuite) TestAcceptRequest() {
}

func (s *TenantSuite) TestRevert() {
}

func TestTenantSuite(t *testing.T) {
	suite.Run(t, &TenantSuite{})
}
