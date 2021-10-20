pub contract FlowContractAudits {

    // Event that is emitted when this contract is created
    pub event ContractInitialized()    
    
    // Event that is emitted when a new Auditor resource is created
    pub event AuditorCreated()

    // Event that is emitted when a new contract audit voucher is created
    pub event AuditVoucherCreated(_ address: Address, codeHash: String)

    // Event that is emitted when a contract audit voucher is removed/used
    pub event AuditVoucherBurned(_ address: Address, codeHash: String)

    // Dictionary of all vouchers currently available
    pub var vouchers: {Address: AuditVoucher}

    // The storage path for the admin resource
    pub let AdminStoragePath: StoragePath

    // The storage Path for auditors' AuditorProxy
    pub let AuditorProxyStoragePath: StoragePath
    
    // The public path for auditors' AuditorProxy capability
    pub let AuditorProxyPublicPath: PublicPath    

    pub struct AuditVoucher {
        pub let codeHash: String
        pub let expiryBlockHeight: UInt64

        init(codeHash: String, expiryBlockHeight: UInt64) {
            self.codeHash = codeHash
            self.expiryBlockHeight = expiryBlockHeight
        }
    }

    pub resource Auditor {
        
        pub fun addAuditVoucher(address: Address, codeHash: String) {

            let voucher = AuditVoucher(codeHash: codeHash, expiryBlockHeight: 1) 

            FlowContractAudits.vouchers.insert(key: address, voucher)

            emit AuditVoucherCreated(address, codeHash: codeHash)
        }

    }

    pub resource interface AuditorProxyPublic {
        pub fun setAuditorCapability(cap: Capability<&Auditor>)
    }    

    pub resource AuditorProxy: AuditorProxyPublic {
        
        access(self) var auditorCapability: Capability<&Auditor>?
        
        pub fun setAuditorCapability(cap: Capability<&Auditor>) {
            self.auditorCapability = cap
        }

        pub fun addAuditVoucher(address: Address, codeHash: String) {
            self.auditorCapability!.borrow()!.addAuditVoucher(address: address, codeHash: codeHash)
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

        pub fun checkAndBurnAuditVoucher(address: Address, codeHash: String): Bool {
            if FlowContractAudits.vouchers[address] != nil {
                FlowContractAudits.vouchers.remove(key: address)                
                return true
            }            
            return false
        }

    }

    init() {
        self.vouchers = {}
        
        self.AdminStoragePath = /storage/contractAuditAdmin
        self.AuditorProxyStoragePath = /storage/contractAuditorProxy
        self.AuditorProxyPublicPath = /public/contractAuditorProxy

        let admin <- create Administrator()
        self.account.save(<-admin, to: self.AdminStoragePath)
        
        emit ContractInitialized()
    }
}