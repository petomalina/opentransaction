package main

import (
	"github.com/petomalina/opentransaction"
	"github.com/petomalina/opentransaction/balance"
)

func main() {
	opentransaction.Serve(&balance.Tenant{})
}
