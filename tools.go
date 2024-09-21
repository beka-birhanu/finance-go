//go: build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	resolver "github.com/beka-birhanu/finance-go/api/graph"
)

type s struct {
	x resolver
}
