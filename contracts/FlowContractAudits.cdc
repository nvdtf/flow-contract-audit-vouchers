pub contract FlowContractAudits {

    // Event that is emitted when this contract is created
    pub event ContractInitialized()    
    
    // Event that is emitted when a new Auditor resource is created
    pub event AuditorCreated()

    // Event that is emitted when a new contract audit voucher is created
    pub event AuditVoucherCreated(_ address: Address, codeHash: String, expiryBlockHeight: UInt64)

    // Event that is emitted when a contract audit voucher is removed/used
    pub event AuditVoucherBurned(_ address: Address, codeHash: String, expiryBlockHeight: UInt64)

    // Dictionary of all vouchers currently available
    pub var vouchers: {String: AuditVoucher}

    // The storage path for the admin resource
    pub let AdminStoragePath: StoragePath

    // The storage Path for auditors' AuditorProxy
    pub let AuditorProxyStoragePath: StoragePath
    
    // The public path for auditors' AuditorProxy capability
    pub let AuditorProxyPublicPath: PublicPath    

    pub struct AuditVoucher {
        pub let address: Address
        pub let codeHash: String
        pub let expiryBlockHeight: UInt64        

        init(address: Address, codeHash: String, expiryBlockHeight: UInt64) {
            self.address = address
            self.codeHash = codeHash
            self.expiryBlockHeight = expiryBlockHeight
        }
    }

    pub fun getVoucherKey(address: Address, codeHash: String) : String {
        return address.toString().concat("-").concat(codeHash)
    }

    pub resource Auditor {
        
        pub fun addAuditVoucher(address: Address, codeHash: String, expiryOffset: UInt64) {

            let expiryBlockHeight = getCurrentBlock().height + expiryOffset

            let key = FlowContractAudits.getVoucherKey(address: address, codeHash: codeHash)
            
            let voucher = AuditVoucher(address: address, codeHash: codeHash, expiryBlockHeight: expiryBlockHeight)            

            FlowContractAudits.vouchers.insert(key: key, voucher)

            emit AuditVoucherCreated(address, codeHash: codeHash, expiryBlockHeight: expiryBlockHeight)
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

        pub fun addAuditVoucher(address: Address, codeHash: String, expiryOffset: UInt64) {
            self.auditorCapability!.borrow()!.addAuditVoucher(address: address, codeHash: codeHash, expiryOffset: expiryOffset)
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
            let key = FlowContractAudits.getVoucherKey(address: address, codeHash: codeHash)
            if FlowContractAudits.vouchers[key] != nil {

                if FlowContractAudits.vouchers[key]!.codeHash == codeHash  {
                    let v = FlowContractAudits.vouchers.remove(key: key)!
                
                    emit AuditVoucherBurned(address, codeHash: v.codeHash, expiryBlockHeight: v.expiryBlockHeight)

                    if getCurrentBlock().height > v.expiryBlockHeight {
                        return false
                    }
                    return true
                }                
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