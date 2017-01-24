package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/http"
	"log"
	"rpcshared"
	"os"
)

var (
	MyName string
	BrokerHost string
)

func init() {
	MyName = Generate(2,"-")
	BrokerHost = os.Getenv("BROKERHOST")
	if len(BrokerHost) == 0 {
		BrokerHost = "trex1:5050"
	}
}

func startServer() {
	be := new(rpcshared.OpenNSFW)
	rpc.Register(be)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":5557")
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	go http.Serve(l, nil)
}

//Start the server, listen forever. 
func main() {
	startServer()
	fmt.Println("[*] Server started. \tMy name:", MyName, "\tBrokerHost: ", BrokerHost)

	meta := make(chan int)
	x := <- meta    /// wait for a while, and listen
	fmt.Println(x)
}
