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
	fmt.Println("Call Router Handle...")

	// 先读取客户端数据, 再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgId(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 创建一个server句柄, 使用zinx的api
	s := znet.NewServer("[zinx v0.5]")

	// 给当前zinx框架注册一个路由
	s.AddRouter(&PingRouter{})

	// 启动服务器
	s.Serve()
}
