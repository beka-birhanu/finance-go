package graph

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
)

func MarshalFloat32(f float32) graphql.Marshaler {
	return graphql.MarshalFloat(float64(f))
}

func UnmarshalFloat32(v interface{}) (float32, error) {
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("%T is not a float64", v)
	}
	return float32(f), nil
}
