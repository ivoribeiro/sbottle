package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	inet "github.com/libp2p/go-libp2p-net"
	libpeer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-protocol"
	"github.com/multiformats/go-multiaddr"
	"os"
)

func handleStream(stream inet.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func main() {
	var peers = map[libpeer.ID]peerstore.PeerInfo{}
	help := flag.Bool("help", false, "Display Help")
	cfg := parseFlags()
	if *help {
		fmt.Printf("Usage: \n   Run './sbottle'\nor Run './sbottle -port [port]'\n")
		os.Exit(0)
	}
	fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.listenHost, cfg.listenPort)
	ctx := context.Background()
	r := rand.Reader
	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}
	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.listenHost, cfg.listenPort))
	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	fmt.Println("This node's multiaddresses:")
	for _, la := range host.Addrs() {
		fmt.Printf(" - %v\n", la)
	}
	fmt.Println()
	if err != nil {
		panic(err)
	}
	// Set a function as stream handler.
	// This function is called when a peer initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), handleStream)
	fmt.Printf("Host ID %s\n", host.ID())
	fmt.Printf("\n[*] Multiaddress: /ip4/%s/tcp/%v/p2p/%s\n", cfg.listenHost, cfg.listenPort, host.ID().Pretty())
	peerChan := initMDNS(ctx, host, cfg.RendezvousString)
	for true {
		// channel from peerChan
		peer := <-peerChan // will block until we discover a peer
		fmt.Println("Found peer:", peer)
		peers[peer.ID] = peer // Add peer to map
		fmt.Println("Peers-> ",peers)
		fmt.Println()
	}
	/*
	if err := host.Connect(ctx, peer); err != nil {
		fmt.Println("Connection failed:", err)
	}
	 //open a stream, this stream will be handled by handleStream other end
	stream, err := host.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))
	if err != nil {
		fmt.Println("Stream open failed", err)
	} else {
	srw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go writeData(rw)
	go readData(rw)
		fmt.Println("Connected to:", peer)
	}
	*/
	select {} //wait here
}
