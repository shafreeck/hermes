package main

import (
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1381")
	_, err = net.ListenTCP("tcp", addr)
	if err != nil {
	}
}
