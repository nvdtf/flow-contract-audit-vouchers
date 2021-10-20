import FlowContractAudits from "../../contracts/FlowContractAudits.cdc"

transaction() {
    
    prepare(adminAccount: AuthAccount) {

        // Create a reference to the admin resource in storage.
        let auditorAdmin = adminAccount.borrow<&FlowContractAudits.Administrator>(from: FlowContractAudits.AdminStoragePath)
            ?? panic("Could not borrow a reference to the admin resource")

        if !auditorAdmin.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test") {
            panic("1")
        }
        if auditorAdmin.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test2") {
            panic("2")
        }
        if auditorAdmin.checkAndBurnAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test") {
            panic("3")
        }

    }

    execute {
    
    }

}