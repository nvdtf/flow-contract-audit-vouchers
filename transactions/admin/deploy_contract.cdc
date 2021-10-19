import FlowContractAudits from "../../contracts/ContractAudits.cdc"

transaction() {
    
    prepare(adminAccount: AuthAccount) {

        if !FlowContractAudits.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test") {
            panic("1")
        }
        if FlowContractAudits.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test2") {
            panic("2")
        }
        if FlowContractAudits.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test") {
            panic("3")
        }

    }

    execute {
    
    }

}