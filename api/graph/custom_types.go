package graph

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalFloat32(f float32) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if _, err := w.Write([]byte(strconv.FormatFloat(float64(f), 'f', -1, 32))); err != nil {
			fmt.Println("Error writing float32:", err)
		}
	})
}

func UnmarshalFloat32(v interface{}) (float32, error) {
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("%T is not a float64", v)
	}
	return float32(f), nil
}
