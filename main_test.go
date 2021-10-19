package main

import (
	"testing"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

func Test(test *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// g.TransactionFromFile("admin/transfer_flow").SignProposeAndPayAsService().UFix64Argument("1.00").AccountArgument("1_auditor").RunPrintEventsFull()

	g.TransactionFromFile("auditor/init").
		SignProposeAndPayAs("1_auditor").
		Test(test).AssertSuccess()

	g.TransactionFromFile("admin/authorize_auditor").
		SignProposeAndPayAsService().
		AccountArgument("1_auditor").
		Test(test).AssertSuccess()

	g.TransactionFromFile("auditor/add_audit").
		SignProposeAndPayAs("1_auditor").
		Test(test).AssertSuccess()

	g.TransactionFromFile("admin/deploy_contract").
		SignProposeAndPayAsService().
		Test(test).AssertSuccess()
}
