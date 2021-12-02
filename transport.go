package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
	}
}

// To construct a simple host with all the default settings, just use `New`
func main() {
	priv, _, err := crypto.GenerateKeyPair(
		crypto.Ed25519, // Select your key type. Ed25519 are nice short
		-1,             // Select key length when possible (i.e. RSA).
	)
	if err != nil {
		panic(err)
	}
	h, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0", // regular tcp connections

		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New))

	if err != nil {
		panic(err)
	}
	//fmt.Println(h2)
	//nickFlag := flag.String("nick", "", "nickname to use in chat. will be generated if empty")
	roomFlag := flag.String("room", "awesome-chat-room", "name of chat room to join")
	flag.Parse()

	ctx := context.Background()

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	fmt.Println("this is pubsub", ps)
	// setup local mDNS discovery
	const DiscoveryInterval = time.Hour

	// DiscoveryServiceTag is used in our mDNS advertisements to discover other chat peers.
	const DiscoveryServiceTag = "pubsub-chat-example"
	if err := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h}).Start(); err != nil {
		panic(err)
	}

	// use the nickname from the cli flag, or a default if blank
	/* nick := *nickFlag
	if len(nick) == 0 {
		nick = defaultNick(h.ID())
	} */

	// join the room from the cli flag, or the flag default
	room := *roomFlag
	fmt.Println(room)
	// join the chat room

	//ctx := context.Background()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

}
