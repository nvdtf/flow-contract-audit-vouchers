package main

import (
	"fmt"
	"testing"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
	"github.com/onflow/cadence"
)

const (
	TestContractCode     = "contract CodyCode {}"
	TestContractCodeSHA3 = "cd1057bd9f593dab406b0a09ffcc7f7468d3ef85021884c4b07430933d94fec0"

	AuditorInitTx                 = "auditor/init"
	AuditorNewAuditTx             = "auditor/new_audit"
	AuditorNewAuditAnyAccountTx   = "auditor/new_audit_any_account"
	AuditorDeleteAuditTx          = "auditor/delete_audit"
	AdminAuthorizeAuditorTx       = "admin/authorize_auditor"
	AdminCleanupExpiredVouchersTx = "admin/cleanup_expired"
	DeveloperDeployContractTx     = "fvm/deploy_contract"

	GetVouchersScript = "get_vouchers"

	AuditorAccount    = "auditor"
	DeveloperAccount  = "developer"
	DeveloperAccount2 = "developer2"
	DeveloperAccount3 = "developer3"

	AuditorCreatedEventName      = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditorCreated"
	AuditVoucherCreatedEventName = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditVoucherCreated"
	AuditVoucherUsedEventName    = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditVoucherUsed"
	AuditVoucherRemovedEventName = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditVoucherRemoved"

	ErrorNoVoucher = "invalid voucher"
)

func TestDeployContract(t *testing.T) {

	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// no voucher on start
	deployAndFail(g, t, DeveloperAccount)

	// init auditor
	authorizeAuditor(g, t)

	// auditor creates new voucher for developer account
	auditContract(g, t, false, false, 10, 19)

	// developer cannot deploy to another account
	deployAndFail(g, t, DeveloperAccount2)
	deployAndFail(g, t, DeveloperAccount3)

	// developer can deploy audited contract
	deploy(g, t, DeveloperAccount, false, 19, false)

	// developer cannot deploy audited contract twice
	deployAndFail(g, t, DeveloperAccount)
}

func TestDeployRecurrentContract(t *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// init auditor
	authorizeAuditor(g, t)

	// auditor adds recurrent voucher for any account
	auditContract(g, t, true, true, 10, 18)

	// developer can deploy audited contract
	deploy(g, t, DeveloperAccount, true, 18, true)

	// developer can deploy audited contract again
	deploy(g, t, DeveloperAccount2, true, 18, true)
	deploy(g, t, DeveloperAccount3, true, 18, true)

	// auditor updates voucher to non-recurrent for any account
	g.TransactionFromFile(AuditorNewAuditAnyAccountTx).
		SignProposeAndPayAs(AuditorAccount).
		StringArgument(TestContractCode).
		BooleanArgument(false).
		UInt64Argument(1).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherCreatedEventName, map[string]interface{}{
			"address":           "",
			"codeHash":          TestContractCodeSHA3,
			"expiryBlockHeight": "13",
			"recurrent":         "false",
		})).
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherRemovedEventName, map[string]interface{}{
			"key":               fmt.Sprintf("any-%s", TestContractCodeSHA3),
			"expiryBlockHeight": "18",
			"recurrent":         "true",
		}))

	// developer deploys and uses voucher
	deploy(g, t, DeveloperAccount, false, 13, true)

	// developer cannot deploy any more
	deployAndFail(g, t, DeveloperAccount)
}

func TestDeleteVoucher(t *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// init auditor
	authorizeAuditor(g, t)

	// auditor adds recurrent voucher
	auditContract(g, t, false, true, 10, 18)

	// developer can deploy audited contract
	deploy(g, t, DeveloperAccount, true, 18, false)

	// delete voucher
	g.TransactionFromFile(AuditorDeleteAuditTx).
		SignProposeAndPayAs(AuditorAccount).
		Argument(cadence.NewOptional(cadence.BytesToAddress(g.Account(DeveloperAccount).Address().Bytes()))).
		StringArgument(TestContractCodeSHA3).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherRemovedEventName, map[string]interface{}{
			"key":               "0x" + g.Account(DeveloperAccount).Address().String() + "-" + TestContractCodeSHA3,
			"expiryBlockHeight": "18",
			"recurrent":         "true",
		}))

	// developer cannot deploy any more
	deployAndFail(g, t, DeveloperAccount)
}

func TestExpiredVouchers(t *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// init auditor
	authorizeAuditor(g, t)

	// auditor adds recurrent voucher for any account
	auditContract(g, t, true, true, 2, 10)

	// developer can deploy audited contract for 2 blocks
	deploy(g, t, DeveloperAccount, true, 10, true)
	deploy(g, t, DeveloperAccount2, true, 10, true)

	// voucher expired
	deployAndFail(g, t, DeveloperAccount3)
}

func TestCleanupExpired(t *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// init auditor
	authorizeAuditor(g, t)

	// auditor adds recurrent voucher for any account
	auditContract(g, t, true, true, 1, 9)

	// check count
	if getVouchersCount(g, t) != 1 {
		t.Fail()
	}

	// cleanup
	g.TransactionFromFile(AdminCleanupExpiredVouchersTx).
		SignProposeAndPayAsService().
		Test(t).
		AssertSuccess()

	// check count, block offset still valid
	if getVouchersCount(g, t) != 1 {
		t.Fail()
	}

	// cleanup
	g.TransactionFromFile(AdminCleanupExpiredVouchersTx).
		SignProposeAndPayAsService().
		Test(t).
		AssertSuccess()

	// verify cleanup
	if getVouchersCount(g, t) != 0 {
		t.Fail()
	}
}

func authorizeAuditor(g *gwtf.GoWithTheFlow, t *testing.T) {
	// auditor init proxy
	g.TransactionFromFile(AuditorInitTx).
		SignProposeAndPayAs(AuditorAccount).
		Test(t).
		AssertSuccess()

	// admin authorizes auditor
	g.TransactionFromFile(AdminAuthorizeAuditorTx).
		SignProposeAndPayAsService().
		AccountArgument(AuditorAccount).
		Test(t).
		AssertSuccess().
		AssertEmitEventName(AuditorCreatedEventName)
}

func deployAndFail(g *gwtf.GoWithTheFlow, t *testing.T, account string) {
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(account).
		StringArgument(TestContractCode).
		Test(t).
		AssertFailure(ErrorNoVoucher)
}

func auditContract(g *gwtf.GoWithTheFlow, t *testing.T, anyAccount bool, recurrent bool, expiryOffset uint64, expiryBlockHeight uint64) {
	var builder gwtf.FlowTransactionBuilder
	var address string

	if anyAccount {
		builder = g.TransactionFromFile(AuditorNewAuditAnyAccountTx)
	} else {
		address = "0x" + g.Account(DeveloperAccount).Address().String()
		builder = g.TransactionFromFile(AuditorNewAuditTx).
			Argument(cadence.NewOptional(cadence.BytesToAddress(g.Account(DeveloperAccount).Address().Bytes())))
	}

	builder.SignProposeAndPayAs(AuditorAccount).
		StringArgument(TestContractCode).
		BooleanArgument(recurrent).
		UInt64Argument(expiryOffset).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherCreatedEventName, map[string]interface{}{
			"address":           address,
			"codeHash":          TestContractCodeSHA3,
			"expiryBlockHeight": fmt.Sprintf("%d", expiryBlockHeight),
			"recurrent":         fmt.Sprintf("%t", recurrent),
		}))
}

func deploy(g *gwtf.GoWithTheFlow, t *testing.T, account string, recurrent bool, expiryBlockHeight uint64, anyAccountVoucher bool) {
	key := fmt.Sprintf("0x%s-%s", g.Account(account).Address().String(), TestContractCodeSHA3)
	if anyAccountVoucher {
		key = fmt.Sprintf("any-%s", TestContractCodeSHA3)
	}
	expiryBlockHeightStr := fmt.Sprintf("%d", expiryBlockHeight)
	recurrentStr := fmt.Sprintf("%t", recurrent)

	result := g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(account).
		StringArgument(TestContractCode).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherUsedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(account).Address().String(),
			"key":               key,
			"expiryBlockHeight": expiryBlockHeightStr,
			"recurrent":         recurrentStr,
		}))

	if !recurrent {
		result.AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherRemovedEventName, map[string]interface{}{
			"key":               key,
			"expiryBlockHeight": expiryBlockHeightStr,
			"recurrent":         recurrentStr,
		}))
	}
}

func getVouchersCount(g *gwtf.GoWithTheFlow, t *testing.T) int {
	countVouchers, err := g.ScriptFromFile(GetVouchersScript).RunReturns()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	c := countVouchers.(cadence.Int)
	return c.Int()
}
