package main

import "TCP-Server/znet"

func main() {
	// 创建一个服务器句柄
	s := znet.NewServer("MMO Game Zinx")

	// 添加连接的创建和销毁的Hook函数

	// 注册一些路由业务

	// 启动服务器
	s.Serve()
}
