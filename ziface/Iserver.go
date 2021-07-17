package ziface

// 定义一个server接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// 路由功能：给当前的服务注册一个路由方法, 供客户端的连接处理使用
	AddRouter(msgId uint32, router Irouter)
	// 获取当前Server的连接管理器
	GetConnMgr() IConnManager

	// 注册 OnConnStart 钩子函数的方法
	SetOnConnStart(func(connection Iconnection))
	// 注册 OnConnStop 钩子函数的方法
	SetOnConnStop(func(connection Iconnection))
	// 调用 OnConnStart 钩子函数的方法
	CallOnConnStart(connection Iconnection)
	// 调用 OnConnStop 钩子函数的方法
	CallOnConnStop(connection Iconnection)
}
