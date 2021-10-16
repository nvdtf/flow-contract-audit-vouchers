package main

import (
	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

func main() {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	g.TransactionFromFile("admin/transfer_flow").SignProposeAndPayAsService().UFix64Argument("1.00").AccountArgument("1_auditor").RunPrintEventsFull()

	// g.TransactionFromFile("auditor/init").SignProposeAndPayAs("auditor").RunPrintEventsFull()
}
