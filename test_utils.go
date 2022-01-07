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
	AuditorDeleteAuditTx          = "auditor/delete_audit"
	AdminAuthorizeAuditorTx       = "admin/authorize_auditor"
	AdminCleanupExpiredVouchersTx = "admin/cleanup_expired"
	DeveloperDeployContractTx     = "fvm/deploy_contract"

	GetVouchersScript = "get_vouchers"

	AuditorAccount    = "auditor"
	DeveloperAccount  = "developer"
	DeveloperAccount2 = "developer2"
	DeveloperAccount3 = "developer3"

	AuditorCreatedEventName = "A.f8d6e0586b0a20c7.FlowContractAudits.AuditorCreated"
	VoucherCreatedEventName = "A.f8d6e0586b0a20c7.FlowContractAudits.VoucherCreated"
	VoucherUsedEventName    = "A.f8d6e0586b0a20c7.FlowContractAudits.VoucherUsed"
	VoucherRemovedEventName = "A.f8d6e0586b0a20c7.FlowContractAudits.VoucherRemoved"

	ErrorNoVoucher = "invalid voucher"
)

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

func auditContract(g *gwtf.GoWithTheFlow, t *testing.T, anyAccount bool, recurrent bool, expiryOffset int, expiryBlockHeight int) {
	builder := g.TransactionFromFile(AuditorNewAuditTx)

	var address string
	if anyAccount {
		builder = builder.Argument(cadence.NewOptional(nil))
	} else {
		address = "0x" + g.Account(DeveloperAccount).Address().String()
		acc := g.Account(DeveloperAccount)
		builder = builder.Argument(cadence.NewOptional(cadence.NewAddress(acc.Address())))
	}

	builder = builder.SignProposeAndPayAs(AuditorAccount).
		StringArgument(TestContractCode).
		BooleanArgument(recurrent)

	var expiryHeight string
	if expiryOffset > 0 {
		expiryHeight = fmt.Sprintf("%d", expiryBlockHeight)
		builder = builder.Argument(cadence.NewOptional(cadence.NewUInt64(uint64(expiryOffset))))
	} else {
		builder = builder.Argument(cadence.NewOptional(nil))
	}

	builder.
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(VoucherCreatedEventName, map[string]interface{}{
			"address":           address,
			"codeHash":          TestContractCodeSHA3,
			"expiryBlockHeight": expiryHeight,
			"recurrent":         fmt.Sprintf("%t", recurrent),
		}))
}

func deploy(g *gwtf.GoWithTheFlow, t *testing.T, account string, recurrent bool, expiryBlockHeight int, anyAccountVoucher bool) {
	key := fmt.Sprintf("0x%s-%s", g.Account(account).Address().String(), TestContractCodeSHA3)
	if anyAccountVoucher {
		key = fmt.Sprintf("any-%s", TestContractCodeSHA3)
	}

	expiryBlockHeightStr := ""
	if expiryBlockHeight > 0 {
		expiryBlockHeightStr = fmt.Sprintf("%d", expiryBlockHeight)
	}

	recurrentStr := fmt.Sprintf("%t", recurrent)

	result := g.TransactionFromFile(DeveloperDeployContractTx).
		SignProposeAndPayAsService().
		AccountArgument(account).
		StringArgument(TestContractCode).
		Test(t).
		AssertSuccess().
		AssertEmitEvent(gwtf.NewTestEvent(VoucherUsedEventName, map[string]interface{}{
			"address":           "0x" + g.Account(account).Address().String(),
			"key":               key,
			"expiryBlockHeight": expiryBlockHeightStr,
			"recurrent":         recurrentStr,
		}))

	if !recurrent {
		result.AssertEmitEvent(gwtf.NewTestEvent(VoucherRemovedEventName, map[string]interface{}{
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
