import FlowContractAudits from "../../contracts/FlowContractAudits.cdc"

transaction() {
    
    let auditorCapability: &FlowContractAudits.AuditorProxy

    prepare(auditorAccount: AuthAccount) {

        self.auditorCapability = auditorAccount
            .borrow<&FlowContractAudits.AuditorProxy>(from: FlowContractAudits.AuditorProxyStoragePath)
            ?? panic("Could not borrow a reference to the admin resource")

    }

    execute {
        self.auditorCapability.addAuditVoucher(address: 0x179b6b1cb6755e31, codeHash: "test", expiryOffset: 1)
    }

}