package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test start...")

	conn, err := net.Dial("tcp4", "127.0.0.1:9190")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}

	for {

		if _, err = conn.Write([]byte("hello zinx v0.2")); err != nil {
			fmt.Println("conn.Write err = ", err)
			continue
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err = ", err)
			continue
		}
		fmt.Printf(" server call back : %s, cnt = %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}
