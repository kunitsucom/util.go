module github.com/kunitsucom/util.go/grpc

go 1.21

require (
	github.com/kunitsucom/util.go v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
)

replace github.com/kunitsucom/util.go => ../../util.go

require github.com/golang/protobuf v1.5.3 // indirect
