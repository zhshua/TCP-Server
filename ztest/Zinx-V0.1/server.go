package main

import (
	"TCP-Server/znet"
)

func main() {
	s := znet.NewServer("[zinx v0.1]")
	s.Serve()
}
