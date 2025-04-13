module github.com/kunitsucom/util.go/x/oauth2

go 1.23.0

toolchain go1.24.1

replace github.com/kunitsucom/util.go => ../../../util.go

require (
	github.com/kunitsucom/util.go v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.29.0
)

require (
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
)
