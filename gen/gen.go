package main

import (
	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-address"
)

func main() {
	if err := gen.WriteTupleEncodersToFile("./cbor_gen.go", "address",
		address.SubnetID{},
	); err != nil {
		panic(err)
	}
}
