package balance

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TenantSuite struct {
	suite.Suite
}

func (s *TenantSuite) TestAccept() {

}

func (s *TenantSuite) TestAcceptRequest() {

}

func (s *TenantSuite) TestRevert() {

}

func TestTenantSuite(t *testing.T) {
	suite.Run(t, &TenantSuite{})
}
