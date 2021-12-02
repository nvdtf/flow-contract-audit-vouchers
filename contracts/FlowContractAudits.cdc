//import Crypto

pub contract FlowContractAudits {

    // Event that is emitted when this contract is created
    pub event ContractInitialized()

    // Event that is emitted when a new Auditor resource is created
    pub event AuditorCreated()

    // Event that is emitted when a new contract audit voucher is created
    pub event AuditVoucherCreated(address: Address?, recurrent: Bool, expiryBlockHeight: UInt64?, codeHash: String)

    // Event that is emitted when a contract audit voucher is used
    pub event AuditVoucherUsed(address: Address, key: String, recurrent: Bool, expiryBlockHeight: UInt64?)

    // Event that is emitted when a contract audit voucher is removed
    pub event AuditVoucherRemoved(key: String, recurrent: Bool, expiryBlockHeight: UInt64?)

    // Dictionary of all vouchers currently available
    access(contract) var vouchers: {String: AuditVoucher}

    // The storage path for the admin resource
    pub let AdminStoragePath: StoragePath

    // The storage Path for auditors' AuditorProxy
    pub let AuditorProxyStoragePath: StoragePath

    // The public path for auditors' AuditorProxy capability
    pub let AuditorProxyPublicPath: PublicPath

    pub struct AuditVoucher {
        pub let address: Address?
        pub let recurrent: Bool
        pub let expiryBlockHeight: UInt64?
        pub let codeHash: String

        init(address: Address?, recurrent: Bool, expiryBlockHeight: UInt64?, codeHash: String) {
            self.address = address
            self.recurrent = recurrent
            self.expiryBlockHeight = expiryBlockHeight
            self.codeHash = codeHash
        }
    }

    pub fun getAllVouchers(): {String: AuditVoucher} {
        return self.vouchers
    }

    pub fun generateVoucherKey(address: Address?, codeHash: String): String {
        if address != nil {
            return address!.toString().concat("-").concat(codeHash)
        }
        return "any-".concat(codeHash)
    }

    pub fun hashContractCode(_ code: String): String {
        return String.encodeHex(HashAlgorithm.SHA3_256.hash(code.utf8))
    }

    pub resource Auditor {

        pub fun addAuditVoucher(address: Address?, recurrent: Bool, expiryOffset: UInt64?, code: String) {

            var expiryBlockHeight: UInt64? = nil
            if expiryOffset != nil {
                expiryBlockHeight = getCurrentBlock().height + expiryOffset!
            }

            let codeHash = FlowContractAudits.hashContractCode(code)

            let key = FlowContractAudits.generateVoucherKey(address: address, codeHash: codeHash)

            let voucher = AuditVoucher(address: address, recurrent: recurrent, expiryBlockHeight: expiryBlockHeight, codeHash: codeHash)

            // TODO update existing voucher if audited with different params -> or should remove first? -> add remove for auditor

            FlowContractAudits.vouchers.insert(key: key, voucher)

            emit AuditVoucherCreated(address: address, recurrent: recurrent, expiryBlockHeight: expiryBlockHeight, codeHash: codeHash)
        }

    }

    pub resource interface AuditorProxyPublic {
        pub fun setAuditorCapability(_ cap: Capability<&Auditor>)
    }

    pub resource AuditorProxy: AuditorProxyPublic {

        access(self) var auditorCapability: Capability<&Auditor>?

        pub fun setAuditorCapability(_ cap: Capability<&Auditor>) {
            self.auditorCapability = cap
        }

        pub fun addAuditVoucher(address: Address?, recurrent: Bool, expiryOffset: UInt64?, code: String) {
            self.auditorCapability!.borrow()!.addAuditVoucher(address: address, recurrent: recurrent, expiryOffset: expiryOffset, code: code)
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
        
        pub fun useVoucherForDeploy(address: Address, code: String): Bool {
            let codeHash = FlowContractAudits.hashContractCode(code)

            var key = FlowContractAudits.generateVoucherKey(address: address, codeHash: codeHash)

            if !FlowContractAudits.vouchers.containsKey(key) {
                key = FlowContractAudits.generateVoucherKey(address: nil, codeHash: codeHash)
                if !FlowContractAudits.vouchers.containsKey(key) {
                    return false
                }
            }

            let v = FlowContractAudits.vouchers[key]!

            if v.codeHash == codeHash  {

                if v.expiryBlockHeight != nil {
                    if getCurrentBlock().height > v.expiryBlockHeight! {
                        FlowContractAudits.vouchers.remove(key: key)
                        emit AuditVoucherRemoved(key: key, recurrent: v.recurrent, expiryBlockHeight: v.expiryBlockHeight)
                        return false
                    }
                }

                if !v.recurrent {
                    FlowContractAudits.vouchers.remove(key: key)
                    emit AuditVoucherRemoved(key: key, recurrent: v.recurrent, expiryBlockHeight: v.expiryBlockHeight)                    
                }
                 
                emit AuditVoucherUsed(address: address, key: key, recurrent: v.recurrent, expiryBlockHeight: v.expiryBlockHeight)
                
                return true
            }

            return false
        }

        pub fun cleanupExpiredVouchers() {
            for key in FlowContractAudits.vouchers.keys {
                let v = FlowContractAudits.vouchers[key]!
                if v.expiryBlockHeight != nil {
                    if getCurrentBlock().height > v.expiryBlockHeight! {
                        FlowContractAudits.vouchers.remove(key: key)
                        emit AuditVoucherRemoved(key: key, recurrent: v.recurrent, expiryBlockHeight: v.expiryBlockHeight)                
                    }
                }
            }
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