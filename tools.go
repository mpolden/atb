// +build tools

//go:generate go run honnef.co/go/tools/cmd/staticcheck -checks inherit ./...

package tools

import (
	_ "honnef.co/go/tools/cmd/staticcheck"
)
