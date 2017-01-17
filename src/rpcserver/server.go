package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/http"
	"log"
	"rpcshared"
)

func startServer() {
	be := new(rpcshared.OpenNSFW)
	rpc.Register(be)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":5556")
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	go http.Serve(l, nil)
}

//Start the server, listen forever. 
func main() {
	startServer()
	meta := make(chan int)
	x := <- meta    /// wait for a while, and listen
	fmt.Println(x)
}
