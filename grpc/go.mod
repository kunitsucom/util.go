module github.com/kunitsucom/util.go/grpc

go 1.21

require (
	github.com/kunitsucom/util.go v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.33.0
)

replace github.com/kunitsucom/util.go => ../../util.go

require golang.org/x/sys v0.17.0 // indirect
