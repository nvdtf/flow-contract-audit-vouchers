pub contract FlowContractAudits {
    pub var vouchers: {Address: String}

    pub let AdminStoragePath: StoragePath

    pub let AuditorStoragePath: StoragePath
    pub let AuditorPublicPath: PublicPath

    pub event ContractAudited(_ address: Address, codeHash: String)

    pub resource Auditor {
        pub fun addAudit(address: Address, codeHash: String) {
            FlowContractAudits.vouchers.insert(key: address, codeHash)
            emit ContractAudited(address, codeHash: codeHash)
        }
    }
    
    pub resource Admin {
        pub fun createAuditor(): @Auditor {
            return <-create Auditor()
        }
    }

    access(account) fun checkAndBurnAuditVoucher(address: Address, codeHash: String): Bool {
        if self.vouchers[address] == codeHash {
            self.vouchers.remove(key: address)
            return true
        }
        return false
    }

    init() {
        self.vouchers = {}  
        self.AdminStoragePath = /storage/auditAdmin
        self.AuditorStoragePath = /storage/auditor
        self.AuditorPublicPath = /public/auditor    
    }
}