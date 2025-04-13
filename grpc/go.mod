module github.com/kunitsucom/util.go/grpc

go 1.23.0

toolchain go1.23.8

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/kunitsucom/util.go v0.0.67
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250409194420-de1ac958c67a
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
)

replace github.com/kunitsucom/util.go => ../../util.go

require (
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250409194420-de1ac958c67a // indirect
)
