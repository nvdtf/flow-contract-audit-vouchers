package main

import (
	"testing"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

func Test(test *testing.T) {
	sampleContractCode := "CodyCode"

	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	g.TransactionFromFile("auditor/init").
		SignProposeAndPayAs("auditor").
		Test(test).AssertSuccess()

	g.TransactionFromFile("admin/authorize_auditor").
		SignProposeAndPayAsService().
		AccountArgument("auditor").
		Test(test).AssertSuccess()

	g.TransactionFromFile("auditor/new_audit").
		SignProposeAndPayAs("auditor").
		AccountArgument("developer").
		StringArgument(sampleContractCode).
		// RunPrintEventsFull()
		Test(test).AssertSuccess()

	g.TransactionFromFile("fvm/deploy_contract").
		SignProposeAndPayAsService().
		AccountArgument("developer").
		StringArgument(sampleContractCode).
		// RunPrintEventsFull()
		Test(test).AssertSuccess()

	// test.Fail()
}
