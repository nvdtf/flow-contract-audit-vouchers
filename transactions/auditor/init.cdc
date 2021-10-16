import FlowContractAudits from "../../contracts/ContractAudits.cdc"

transaction {

    prepare(auditor: AuthAccount) {

        let auditorProxy <- FlowContractAudits.createAuditorProxy()

        auditor.save(
            <- auditorProxy, 
            to: FlowContractAudits.ContractAuditorProxyStoragePath,
        )
            
        auditor.link<&FlowContractAudits.AuditorProxy{FlowContractAudits.AuditorProxyPublic}>(
            FlowContractAudits.ContractAuditorProxyPublicPath,
            target: FlowContractAudits.ContractAuditorProxyStoragePath
        )
    }
}