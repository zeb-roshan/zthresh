package main

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"sync"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
)

//var fmt = log.fmt("rendezvous")

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")
	fmt.Println("Connect with ", stream.ID())
	// Create a buffer stream for non blocking read and write.
	//rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	/*
		go readData(rw, stream)
		go writeData(rw, stream)
	*/
	// 'stream' will stay open until you close it (or the other side closes it).
}

func main() {

	//input threshold, name and user
	//fmt.Println("Enter threshold:")
	/* var input_peers int
	fmt.Scanln(&input_peers)
	for input_peers < 2 {
		fmt.Println("Threshold cannot be less than 2, Please re-enter")
		fmt.Scanln(&input_peers)
	} */
	fmt.Println(rand.Intn(9999-1000) + 1000)

	port := strconv.Itoa(rand.Intn(9999-1000) + 1000)
	//fmt.Println("Port: ", reflect.TypeOf(port))

	//TCP connection
	priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	h, _ := libp2p.New(
		//context.Background(),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/"+string(port)),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Identity(priv),
	)

	ctx := context.Background()
	h.SetStreamHandler(protocol.ID("tss/1"), handleStream)
	kademliaDHT, err := dht.New(cont(), h)
	if err != nil {
		panic(err)
	}
	fmt.Println("type is ", h.Addrs(), h.ID().String())
	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	fmt.Println("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				fmt.Println(err)
			} else {
				//fmt.Println("Connection established with bootstrap node:", *peerinfo)
				return
			}
		}()
	}
	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	fmt.Println("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(cont(), routingDiscovery, "base")
	fmt.Println("Successfully announced!")

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	fmt.Println("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(cont(), "base")
	if err != nil {
		panic(err)
	}

	for peer := range peerChan {

		if peer.ID == h.ID() {
			fmt.Println("And here ", reflect.TypeOf(h.ID), reflect.TypeOf(peer.ID))

			continue
		}

		_, err := h.NewStream(cont(), peer.ID, protocol.ID("tss/1"))

		if err != nil {
			//fmt.Println("Connection failed:", err)
			continue
		} else {
			fmt.Println("Found ", peer)
			fmt.Println("Connecting to:", peer.Addrs[0].String()+"/p2p/"+peer.ID.String())
			/* rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			fmt.Println(rw) */
			peerlist = append(peerlist, peers{id: len(peerlist) + 1, addr: peer.Addrs[0].String() + "/p2p/" + peer.ID.String(), connect: false})
			peerattempt = append(peerattempt, peer.ID.String())

			go attemptconnection(len(peerlist) + 1)
			/*
				go writeData(rw, stream)
				go readData(rw, stream) */
		}

		//fmt.Println("Connected to:", peer)
	}

	go handle_write()
	select {}
}

func cont() context.Context {

	return context.Background()
}

//var peers []string
type peers struct {
	id      int
	addr    string
	connect bool
}

var peerlist []peers

var peerattempt []string

func handle_write() {
	for {
		time.Sleep(10 * time.Second)
		fmt.Println(peerlist)
	}
}

func attemptconnection(id int) {

	fmt.Println("Attempting connection")

	//pass code entry

	//Connect and show proof - Add keccak auth?

}
