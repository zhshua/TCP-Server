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

func (this *PingRouter) PreHandle(request ziface.Irequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping PreHandle\n"))
	if err != nil {
		fmt.Println("call back before handle error", err)
	}
}

func (this *PingRouter) Handle(request ziface.Irequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping Handle\n"))
	if err != nil {
		fmt.Println("call back ping handle error", err)
	}
}

func (this *PingRouter) PostHandle(request ziface.Irequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping PostHandle\n"))
	if err != nil {
		fmt.Println("call back post handle error", err)
	}
}

func main() {
	// 创建一个server句柄, 使用zinx的api
	s := znet.NewServer("[zinx v0.3]")

	// 给当前zinx框架注册一个路由
	s.AddRouter(&PingRouter{})

	// 启动服务器
	s.Serve()
}
