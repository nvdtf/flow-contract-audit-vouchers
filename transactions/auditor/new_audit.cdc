import FlowContractAudits from "../../contracts/FlowContractAudits.cdc"

transaction(address: Address, code: String) {
    
    let auditorCapability: &FlowContractAudits.AuditorProxy

    prepare(auditorAccount: AuthAccount) {

        self.auditorCapability = auditorAccount
            .borrow<&FlowContractAudits.AuditorProxy>(from: FlowContractAudits.AuditorProxyStoragePath)
            ?? panic("Could not borrow a reference to the admin resource")

    }

    execute {
        self.auditorCapability.addAuditVoucher(address: address, recurrent: false, expiryOffset: 1, code: code)        
    }

}