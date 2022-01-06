import FlowContractAudits from "../../contracts/FlowContractAudits.cdc"

transaction() {
    let auditorAdmin: &FlowContractAudits.Administrator
    
    prepare(adminAccount: AuthAccount) {
        // Create a reference to the admin resource in storage.
        self.auditorAdmin = adminAccount.borrow<&FlowContractAudits.Administrator>(from: FlowContractAudits.AdminStoragePath)
            ?? panic("Could not borrow a reference to the admin resource")        
    }

    execute {
        self.auditorAdmin.cleanupExpiredVouchers()             
    }
}