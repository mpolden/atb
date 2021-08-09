// +build tools

// Disable SA5008 because cmd packages has a duplicate "choice" tag
//go:generate go run honnef.co/go/tools/cmd/staticcheck -checks inherit,-SA5008 ./...

package tools

import (
	_ "honnef.co/go/tools/cmd/staticcheck"
)
