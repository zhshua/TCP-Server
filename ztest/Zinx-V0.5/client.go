package main

import (
	"TCP-Server/znet"
	"fmt"
	"io"
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

		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx-V0.5 Client")))
		if err != nil {
			fmt.Println("Pack error: ", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err = io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msg error", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.MsgData = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.MsgData); err != nil {
				fmt.Println("read msg data error", err)
				return
			}

			fmt.Printf("----> Recv Server Msg : ID = %d, len = %d, data = %s\n", msg.Id, msg.DataLen, string(msg.MsgData))

		}

		time.Sleep(1 * time.Second)
	}
}
