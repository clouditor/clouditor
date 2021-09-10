package evidenceStore

//go:generate protoc -I ../../proto -I ../../third_party evidence_store.proto --go_out=../.. --go-grpc_out=../.. --go_opt=Mevidence.proto=clouditor.io/clouditor/api/assessment --go-grpc_opt=Mevidence.proto=clouditor.io/clouditor/api/assessment --openapi_out=../../openapi/evidenceStore
