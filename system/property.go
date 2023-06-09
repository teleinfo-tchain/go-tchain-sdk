package system

const (
	AllianceContract     = "did:bid:sf29zfRvb9hFmcD2A51HRHjFKZrqtx78p"
	ElectionContract     = "did:bid:sfMw1S8VY6eVyccpgKQphBkpg9BM7GF6"
	CertificateContract  = "did:bid:sf4Maf18h4Y5jDw94PRd4gEHwfPbANhk"
	DocumentContract     = "did:bid:sfcR79cThpmU8pjgvbxXZUsh4jTVoLGz"
	SensitiveContract    = "did:bid:sf2BJooAfaZoUiCHhtwxP4pUJBA19u5Fh"
	SuperManagerContract = "did:bid:sfmQGYcy2knexz9iJxd2fmuQa3GQugK8"
	SubChainContract     = "did:bid:sfiGyvAxK7d2PT8UDpEssL6MVdj6F8g1"
)

const AllianceAbiJSON = `[
{"constant": false,"name":"registerDirector","inputs":[{"name":"id","type":"string"},{"name":"publicKey","type":"string"},{"name":"companyName","type":"string"},{"name":"companyCode","type":"string"}],"outputs":[],"type":"function"}, 
{"constant": false,"name":"upgradeDirector","inputs":[{"name":"director","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revoke","inputs":[{"name":"member","type":"string"}, {"name":"revokeReason","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"setWeights","inputs":[{"name":"directorWeights","type":"uint64"},{"name":"viceWeights","type":"uint64"},{"name":"directorGeneralWeights","type":"uint64"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"allianceEvent","type":"event"}
]`

const ElectionAbiJSON = `[
{"constant": false,"name":"registerTrustNode","inputs":[{"name":"id","type":"string"},{"name":"apply","type":"string"},{"name":"publicKey","type":"string"},{"name":"nodeName","type":"string"},{"name":"messageSha3","type":"string"},{"name":"signature","type":"string"},{"name":"website","type":"string"},{"name":"nodeType","type":"uint64"},{"name":"companyName","type":"string"},{"name":"companyCode","type":"string"},{"name":"ip","type":"string"},{"name":"port","type":"uint64"},{"name":"ipType","type":"uint64"}],"outputs":[],"type":"function"}, 
{"constant": false,"name":"deleteTrustNode","inputs":[{"name":"trustNode","type":"string"}, {"name":"revokeReason","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"applyCandidate","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"cancelCandidate","inputs":[{"name":"candidate","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"voteCandidate","inputs":[{"name":"candidates","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"cancelConsensusNode","inputs":[{"name":"consensus","type":"string"}, {"name":"cancelConsensusReason","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"setDeadline","inputs":[{"name":"deadline","type":"uint64"}],"outputs":[],"type":"function"},
{"constant": false,"name":"extractOwnBounty","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"electEvent","type":"event"},
{"constant": false,"name":"issueAdditionalBounty","inputs":[],"outputs":[],"type":"function"}
]`

const CertificateAbiJSON = `[
{"constant": false,"name":"issueCertificate","inputs":[{"name":"id","type":"string"},{"name":"context","type":"string"},{"name":"subject","type":"string"},{"name":"period","type":"uint64"},{"name":"issuer_algorithm","type":"string"},{"name":"issuer_signature","type":"string"},{"name":"subject_public_key","type":"string"},{"name":"subject_algorithm","type":"string"},{"name":"subject_signature","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificate","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revokedCertificates","inputs":[],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"cerdEvent","type":"event"}
]`

const DocAbiJSON = `[
{"constant": false,"name":"init","inputs":[{"name":"bid_type","type":"uint64"}],"outputs":[],"type":"function"},
{"constant": false,"name":"setBidName","inputs":[{"name":"id","type":"string"},{"name":"bid_name","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addPublic","inputs":[{"name":"id","type":"string"},{"name":"public_type","type":"string"},{"name":"public_auth","type":"string"},{"name":"public_key","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delPublic","inputs":[{"name":"id","type":"string"},{"name":"public_key","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addAuth","inputs":[{"name":"id","type":"string"},{"name":"auth","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delAuth","inputs":[{"name":"id","type":"string"},{"name":"auth","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addService","inputs":[{"name":"id","type":"string"},{"name":"service_id","type":"string"},{"name":"service_type","type":"string"},{"name":"service_endpoint","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delService","inputs":[{"name":"id","type":"string"},{"name":"service_id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addProof","inputs":[{"name":"id","type":"string"},{"name":"proof_type","type":"string"},{"name":"proof_creator","type":"string"},{"name":"proof_sign","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delProof","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"addExtra","inputs":[{"name":"id","type":"string"},{"name":"extra","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delExtra","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"enable","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"disable","inputs":[{"name":"id","type":"string"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"bidEvent","type":"event"}
]`

const SensitiveWordsAbiJSON = `[
{"constant": false,"name":"addWords","inputs":[{"name":"word","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"delWord","inputs":[{"name":"word","type":"string"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"sensitiveEvent","type":"event"}
]`

const ManagerAbiJSON = `[
{"constant": false,"name":"enable","inputs":[{"name":"contract_address","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"disable","inputs":[{"name":"contract_address","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"power","inputs":[{"name":"user_address","type":"string"},{"name":"power","type":"uint64"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"superManagerEvent","type":"event"}
]`

const SubChainAbiJSON = `[
{"constant": false,"name":"applySubChain","inputs":[{"name":"id","type":"string"},{"name":"apply","type":"string"},{"name":"subChainName","type":"string"},{"name":"chainCode","type":"string"},{"name":"chainIndustry","type":"string"},{"name":"chainFramework","type":"string"},{"name":"consensus","type":"string"},{"name":"chainMsgHash","type":"string"}],"outputs":[],"type":"function"}, 
{"constant": false,"name":"voteSubChain","inputs":[{"name":"candidates","type":"string"}],"outputs":[],"type":"function"},
{"constant": false,"name":"setDeadline","inputs":[{"name":"deadline","type":"uint64"}],"outputs":[],"type":"function"},
{"constant": false,"name":"revoke","inputs":[{"name":"subChainId","type":"string"}, {"name":"revokeReason","type":"string"}],"outputs":[],"type":"function"},
{"anonymous":false,"inputs":[{"indexed":false,"name":"method_name","type":"string"},{"indexed":false,"name":"status","type":"uint32"},{"indexed":false,"name":"reason","type":"string"}],"name":"subChainEvent","type":"event"}
]`
