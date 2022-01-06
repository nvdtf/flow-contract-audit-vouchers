import FlowContractAudits from "../../contracts/FlowContractAudits.cdc"

transaction(address: Address?, codeHash: String) {
    
    let auditorCapability: &FlowContractAudits.AuditorProxy

    prepare(auditorAccount: AuthAccount) {
        self.auditorCapability = auditorAccount
            .borrow<&FlowContractAudits.AuditorProxy>(from: FlowContractAudits.AuditorProxyStoragePath)
            ?? panic("Could not borrow a reference to the admin resource")
    }

    execute {
        self.auditorCapability.deleteVoucher(address: address, codeHash: codeHash)        
    }
}
 