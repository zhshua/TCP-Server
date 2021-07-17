package main

import (
	"TCP-Server/ziface"
	"TCP-Server/znet"
	"fmt"
)

// 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.Irequest) {
	fmt.Println("Call PingRouter Handle...")

	// 先读取客户端数据, 再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgId(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.Irequest) {
	fmt.Println("Call HelloRouter Handle...")

	// 先读取客户端数据, 再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgId(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("hello...hello...hello"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnBegin(conn ziface.Iconnection) {
	fmt.Println("-----> DoConnBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnBegin")); err != nil {
		fmt.Println(err)
	}
}

func DoConnEnd(conn ziface.Iconnection) {
	fmt.Println("-----> DoConnEnd is Called...")
	fmt.Println("ConnID = ", conn.GetConnID(), "is Offline...")
}

func main() {
	// 创建一个server句柄, 使用zinx的api
	s := znet.NewServer("[zinx v0.9]")

	// 注册连接Hook函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnEnd)

	// 给当前zinx框架注册路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 启动服务器
	s.Serve()
}
