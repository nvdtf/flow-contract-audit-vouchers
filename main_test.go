package main

import (
	"testing"

	"github.com/bjartek/go-with-the-flow/v2/gwtf"
)

const (
	TestContractCode     = "contract CodyCode {}"
	TestContractCodeSHA3 = "cd1057bd9f593dab406b0a09ffcc7f7468d3ef85021884c4b07430933d94fec0"

	AuditorInitTx             = "auditor/init"
	AuditorNewAuditTx         = "auditor/new_audit"
	AdminAuthorizeAuditorTx   = "admin/authorize_auditor"
	DeveloperDeployContractTx = "fvm/deploy_contract"

	AuditorAccount   = "auditor"
	DeveloperAccount = "developer"

	AuditorCreatedEventName      = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditorCreated"
	AuditVoucherCreatedEventName = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditVoucherCreated"
	AuditVoucherUsedEventName    = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditVoucherUsed"
	AuditVoucherRemovedEventName = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditVoucherRemoved"

	ErrorNoVoucher = "invalid voucher"
)

func TestDeployContract(t *testing.T) {

	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// no voucher on start
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).
		AssertFailure(ErrorNoVoucher)

	// init auditor
	authorizeAuditor(g, t)

	// auditor creates new voucher
	g.TransactionFromFile(AuditorNewAuditTx).
		SignProposeAndPayAs(AuditorAccount).
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		BooleanArgument(false).
		UInt64Argument(1).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherCreatedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(DeveloperAccount).Address().String(),
			"codeHash":          TestContractCodeSHA3,
			"expiryBlockHeight": "8",
			"recurrent":         "false",
		}))

	// developer can deploy audited contract
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherRemovedEventName, map[string]interface{}{
			"key":               "0x" + g.Account(DeveloperAccount).Address().String() + "-" + TestContractCodeSHA3,
			"expiryBlockHeight": "8",
			"recurrent":         "false",
		})).
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherUsedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(DeveloperAccount).Address().String(),
			"key":               "0x" + g.Account(DeveloperAccount).Address().String() + "-" + TestContractCodeSHA3,
			"expiryBlockHeight": "8",
			"recurrent":         "false",
		}))

	// developer cannot deploy audited contract twice
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).
		AssertFailure(ErrorNoVoucher)
}

func TestDeployRecurrentContract(t *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	// init auditor
	authorizeAuditor(g, t)

	// auditor adds recurrent voucher
	g.TransactionFromFile(AuditorNewAuditTx).
		SignProposeAndPayAs(AuditorAccount).
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		BooleanArgument(true).
		UInt64Argument(10).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherCreatedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(DeveloperAccount).Address().String(),
			"codeHash":          TestContractCodeSHA3,
			"expiryBlockHeight": "16",
			"recurrent":         "true",
		}))

	// developer can deploy audited contract
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherUsedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(DeveloperAccount).Address().String(),
			"key":               "0x" + g.Account(DeveloperAccount).Address().String() + "-" + TestContractCodeSHA3,
			"expiryBlockHeight": "16",
			"recurrent":         "true",
		}))

	// developer can deploy audited contract again
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherUsedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(DeveloperAccount).Address().String(),
			"key":               "0x" + g.Account(DeveloperAccount).Address().String() + "-" + TestContractCodeSHA3,
			"expiryBlockHeight": "16",
			"recurrent":         "true",
		}))

	// auditor updates voucher to non-recurrent
	g.TransactionFromFile(AuditorNewAuditTx).
		SignProposeAndPayAs(AuditorAccount).
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		BooleanArgument(false).
		UInt64Argument(1).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherCreatedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(DeveloperAccount).Address().String(),
			"codeHash":          TestContractCodeSHA3,
			"expiryBlockHeight": "10",
			"recurrent":         "false",
		})).
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherRemovedEventName, map[string]interface{}{
			"key":               "0x" + g.Account(DeveloperAccount).Address().String() + "-" + TestContractCodeSHA3,
			"expiryBlockHeight": "16",
			"recurrent":         "true",
		}))

	// developer deploys and uses voucher
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherRemovedEventName, map[string]interface{}{
			"key":               "0x" + g.Account(DeveloperAccount).Address().String() + "-" + TestContractCodeSHA3,
			"expiryBlockHeight": "10",
			"recurrent":         "false",
		})).
		AssertEmitEvent(gwtf.NewTestEvent(AuditVoucherUsedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(DeveloperAccount).Address().String(),
			"key":               "0x" + g.Account(DeveloperAccount).Address().String() + "-" + TestContractCodeSHA3,
			"expiryBlockHeight": "10",
			"recurrent":         "false",
		}))

	// developer cannot deploy any more
	g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(DeveloperAccount).
		StringArgument(TestContractCode).
		Test(t).
		AssertFailure(ErrorNoVoucher)

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
