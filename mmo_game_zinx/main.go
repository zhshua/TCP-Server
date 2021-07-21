package main

import (
	"TCP-Server/mmo_game_zinx/apis"
	"TCP-Server/mmo_game_zinx/core"
	"TCP-Server/ziface"
	"TCP-Server/znet"
	"fmt"
)

// 注册创建连接之后的Hook函数
func OnConnectionAdd(conn ziface.Iconnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)

	// 给客户端发送MsgId:1的消息: 同步当前Player的id给客户端
	player.SyncPid()

	// 给客户端发送MsgId:200的消息: 同步当前Player的初始位置给客户端
	player.BroadCastStartPosition()

	// 将当前新上线的玩家添加到在线玩家集合
	core.WorldMgrObj.AddPlayer(player)

	// 设置连接额外的属性, 将玩家的id绑定在该连接上
	conn.SetProperty("pid", player.Pid)

	// 同步自己的位置信息给周围九宫格内的玩家
	player.SyncSurrounding()

	fmt.Println("---------> Player pid = ", player.Pid, "is online <---------")
}

// 注册销毁连接之前的Hook函数
func OnConnectionLost(conn ziface.Iconnection) {
	// 通过连接属性拿到pid
	pid, _ := conn.GetProperty("pid")
	// 通过pid得到对应的player对象
	player := core.WorldMgrObj.GetOnlinePlayerByPid(pid.(int32))

	// 触发玩家下线业务
	player.Offline()

	fmt.Println("---------> Player pid = ", player.Pid, "is offline <---------")
}

func main() {
	// 创建一个服务器句柄
	s := znet.NewServer("MMO Game Zinx")

	// 添加连接的创建和销毁的Hook函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)

	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	// 启动服务器
	s.Serve()
}
