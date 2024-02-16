//go:build tools
// +build tools

// This file uses the recommended method for tracking developer tools in a Go
// module.
//
// REF: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/mgechev/revive"
<<<<<<< HEAD
=======
	_ "golang.org/x/tools/cmd/goimports"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
>>>>>>> 8649e17 (feat: Upgrade tia branch to cosmos sdk v0.50.1 (#382))
	_ "mvdan.cc/gofumpt"
)
