package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
)

func main() {

	// The context governs the lifetime of the libp2p node
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// To construct a simple host with all the default settings, just use `New`
	h, err := libp2p.New(ctx,
		// Set your own listen address
		// The config takes an array of addresses, specify as many as you want.
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9000"), )
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello World, my hosts ID is %s\n", h.ID())
}
