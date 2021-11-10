package main

import (
	"testing"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

const (
	TestContractCode = "CodyCode"

	AuditorInitTx             = "auditor/init"
	AuditorNewAuditTx         = "auditor/new_audit"
	AdminAuthorizeAuditorTx   = "admin/authorize_auditor"
	DeveloperDeployContractTx = "fvm/deploy_contract"

	AuditorAccount   = "auditor"
	DeveloperAccount = "developer"

	ErrorNoVoucher = "invalid voucher"
)

func TestDeployContract(t *testing.T) {

	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// no voucher on start
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).AssertFailure(ErrorNoVoucher)

	// auditor init proxy
	g.TransactionFromFile(AuditorInitTx).
		SignProposeAndPayAs(AuditorAccount).
		Test(t).AssertSuccess()

	// admin authorizes auditor
	g.TransactionFromFile(AdminAuthorizeAuditorTx).
		SignProposeAndPayAsService().
		AccountArgument(AuditorAccount).
		Test(t).AssertSuccess()

	// auditor creates new voucher
	g.TransactionFromFile(AuditorNewAuditTx).
		SignProposeAndPayAs(AuditorAccount).
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).AssertSuccess()

	// developer can deploy audited contract
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).AssertSuccess()

	// developer cannot deploy audited contract twice
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).AssertFailure(ErrorNoVoucher)
}
