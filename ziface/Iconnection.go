package ziface

import "net"

type Iconnection interface {
	// 启动连接, 开始当前连接的工作
	Start()

	// 停止连接, 结束当前连接的工作
	Stop()

	// 获取当前连接所绑定的套接字信息
	GetTCPConnection() *net.TCPConn

	// 获取当前连接模块的连接ID
	GetConnID() uint32

	// 获取客户端的TCP端口和地址
	RemoteAddr() net.Addr

	// 发送数据, 将数据发送给远程的客户端
	Send(data []byte) error
}

// 定义一个当前连接对应的业务处理函数
type HandleFunc func(*net.TCPConn, []byte, int) error
