import FlowContractAudits from "../../contracts/ContractAudits.cdc"

transaction(auditorAddress: Address) {

    let resourceStoragePath: StoragePath
    let capabilityPrivatePath: CapabilityPath
    let auditorCapability: Capability<&FlowContractAudits.Auditor>

    prepare(adminAccount: AuthAccount) {

        // These paths must be unique within the FUSD contract account's storage
        self.resourceStoragePath = /storage/auditor     // e.g. /storage/minter_01
        self.capabilityPrivatePath = /private/auditor // e.g. private/minter_01

        // Create a reference to the admin resource in storage.
        let auditorAdmin = adminAccount.borrow<&FlowContractAudits.Administrator>(from: FlowContractAudits.AdminStoragePath)
            ?? panic("Could not borrow a reference to the admin resource")

        // Create a new minter resource and a private link to a capability for it in the admin's storage.
        let auditor <- auditorAdmin.createNewAuditor()
        adminAccount.save(<- auditor, to: self.resourceStoragePath)
        self.auditorCapability = adminAccount.link<&FlowContractAudits.Auditor>(
            self.capabilityPrivatePath,
            target: self.resourceStoragePath
        ) ?? panic("Could not link auditor")

    }

    execute {
        // This is the account that the capability will be given to
        let auditorAccount = getAccount(auditorAddress)

        let capabilityReceiver = auditorAccount.getCapability
            <&FlowContractAudits.AuditorProxy{FlowContractAudits.AuditorProxyPublic}>
            (FlowContractAudits.ContractAuditorProxyPublicPath)!
            .borrow() ?? panic("Could not borrow capability receiver reference")

        capabilityReceiver.setAuditorCapability(cap: self.auditorCapability)
    }

}