package main

import (
	"fmt"
	"time"
	"net"
	"net/rpc"
	"net/http"
	"log"
	"rpcshared"
	"os"

	// Log results to our webserver
	"rpclogger"
)

var (
	MyName string
	BrokerHost string
	MyType	string
)

func init() {
	MyName = Generate(2,"-")
	BrokerHost = os.Getenv("BROKERHOST")
	if len(BrokerHost) == 0 {
		BrokerHost = "trex1:5050"
	}
	MyType = "OpenNSFW Score"
}

func startServer() {
	on := new(rpcshared.OpenNSFW)
	rpc.Register(on)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":5557")
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	go http.Serve(l, nil)
	go PeriodicUpdate(on)
}

// Send the whole request history periodically 
// TODO: Decay the RequestHistory buffer. This struct will eventually get huge..
func PeriodicUpdate(myRPCInstance *rpcshared.OpenNSFW) {
	for {
		time.Sleep(time.Millisecond * 5000)
		rpclogger.SubmitReport(BrokerHost, MyName, MyType, myRPCInstance.RequestHistory)
	}
}

//Start the server, listen forever. 
func main() {
	startServer()
	fmt.Println("[*] Server started. \tMy name:", MyName, "\tBrokerHost: ", BrokerHost)

	meta := make(chan int)
	x := <- meta    /// wait for a while, and listen
	fmt.Println(x)
}
