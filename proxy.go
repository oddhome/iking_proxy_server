package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
)

var localAddress *string = flag.String("l", "localhost:4321", "Local address")
var remoteAddress *string = flag.String("r", "203.70.8.215:4321", "Remote address")

func main() {
	flag.Parse()

	fmt.Printf("Listening: %v\nProxying %v\n", *localAddress, *remoteAddress)

	addr, err := net.ResolveTCPAddr("tcp", *localAddress)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}

		go proxyConnection(conn)
	}

}

func proxyConnection(conn *net.TCPConn) {
	rAddr, err := net.ResolveTCPAddr("tcp", *remoteAddress)
	if err != nil {
		panic(err)
	}

	rConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		panic(err)
	}

	defer rConn.Close()

	// Request loop
	go func() {
		for {
			data := make([]byte, 1024*1024)
			n, err := conn.Read(data)
			if err != nil {
				panic(err)
			}
			rConn.Write(data[:n])
			log.Printf("sent:\n%v", hex.Dump(data[:n]))
		}
	}()

	// Response loop
	for {
		data := make([]byte, 1024*1024)
		n, err := rConn.Read(data)
		if err != nil {
			panic(err)
		}
		//Receive Match and Replace


		//Display RAW Data
		conn.Write(data[:n])
		log.Printf("received:\n%v", hex.Dump(data[:n]))
	}

}