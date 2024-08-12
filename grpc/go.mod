module github.com/kunitsucom/util.go/grpc

go 1.21

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0
	github.com/kunitsucom/util.go v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240723171418-e6d459c13d2a
	google.golang.org/grpc v1.64.1
	google.golang.org/protobuf v1.34.2
)

replace github.com/kunitsucom/util.go => ../../util.go

require (
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240723171418-e6d459c13d2a // indirect
)
