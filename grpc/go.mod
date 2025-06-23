module github.com/kunitsucom/util.go/grpc

go 1.23.0

toolchain go1.23.8

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.0
	github.com/kunitsucom/util.go v0.0.67
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

replace github.com/kunitsucom/util.go => ../../util.go

require (
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822 // indirect
)
