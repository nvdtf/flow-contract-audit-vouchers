pub contract FlowContractAudits {
    pub var vouchers: {Address: String}

    pub let AdminStoragePath: StoragePath

    pub let ContractAuditorStoragePath: StoragePath
    pub let ContractAuditorPublicPath: PublicPath

    pub event ContractInitialized()
    pub event AuditorCreated()
    pub event ContractAudited(_ address: Address, codeHash: String)    

    pub resource Auditor {
        pub fun addAuditVoucher(address: Address, codeHash: String) {
            FlowContractAudits.vouchers.insert(key: address, codeHash)
            emit ContractAudited(address, codeHash: codeHash)
        }
    }  

    access(account) fun checkAndBurnAuditVoucher(address: Address, codeHash: String): Bool {
        if self.vouchers[address] == codeHash {
            self.vouchers.remove(key: address)
            return true
        }
        return false
    }

    pub resource interface AuditorProxyPublic {
        pub fun setAuditorCapability(cap: Capability<&Auditor>)
    }    

    pub resource AuditorProxy: AuditorProxyPublic {
        
        access(self) var auditorCapability: Capability<&Auditor>?
        
        pub fun setAuditorCapability(auditorCap: Capability<&Auditor>) {
            self.auditorCapability = auditorCap
        }

        pub fun addAuditVoucher(address: Address, codeHash: String) {
            return <- self.auditorCapability!
            .borrow()!
            .addAuditVoucher(address: address, codeHash: codeHash)
        }

        init() {
            self.auditorCapability = nil
        }

    }

    pub fun createAuditorProxy(): @AuditorProxy {
        return <- create AuditorProxy()
    }
    
    pub resource Administrator {
        pub fun createNewAuditor(): @Auditor {
            emit AuditorCreated()
            return <-create Auditor()
        }
    }

    init() {
        self.vouchers = {}  
        self.AdminStoragePath = /storage/contractAuditAdmin
        self.ContractAuditorProxyStoragePath = /storage/contractAuditorProxy
        self.ContractAuditorProxyPublicPath = /public/contractAuditorProxy

        let admin <- create Administrator()
        self.account.save(<-admin, to: self.AdminStoragePath)
        
        emit ContractInitialized()
    }
}