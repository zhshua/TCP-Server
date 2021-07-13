package znet

import (
	"TCP-Server/ziface"
	"fmt"
	"net"
)

type Connection struct {
	// 当前连接的套接字
	Conn *net.TCPConn

	// 当前连接的ID
	ConnID uint32

	// 当前连接是否关闭
	isClosed bool

	// 当前连接对应的业务处理函数
	handleFunc ziface.HandleFunc

	// 告知当前连接已经退出的channel, 连接如果要退出的话,通过管道告知一下
	ExitChan chan bool
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		handleFunc: callback_api,
		ExitChan:   make(chan bool),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Groutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("c.Conn.Read err = ", err)
			continue
		}

		// 读取完数据后, 调用当前连接所绑定的业务处理API, 让他去处理他的业务
		if err = c.handleFunc(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID ", c.ConnID, "handle err = ", err)
			break
		}
	}
}

// 启动连接, 开始当前连接的工作
func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID", c.Conn)
	go c.StartReader()
}

// 停止连接, 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果连接已经是关闭的
	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()
}

// 获取当前连接所绑定的套接字信息
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取客户端的TCP端口和地址
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据, 将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
