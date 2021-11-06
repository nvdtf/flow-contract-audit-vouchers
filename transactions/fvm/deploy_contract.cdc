import FlowContractAudits from "../../contracts/FlowContractAudits.cdc"

transaction(address: Address, code: String) {

    let auditorAdmin: &FlowContractAudits.Administrator
    
    prepare(adminAccount: AuthAccount) {

        // Create a reference to the admin resource in storage.
        self.auditorAdmin = adminAccount.borrow<&FlowContractAudits.Administrator>(from: FlowContractAudits.AdminStoragePath)
            ?? panic("Could not borrow a reference to the admin resource")

        // if !auditorAdmin.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test") {
        //     panic("1")
        // }
        // if auditorAdmin.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test2") {
        //     panic("2")
        // }
        // if auditorAdmin.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test") {
        //     panic("3")
        // }

    }

    execute {
        if !self.auditorAdmin.checkAndBurnAuditVoucher(address: address, code: code) {
            panic("invalid voucher")
        }    
    }

}