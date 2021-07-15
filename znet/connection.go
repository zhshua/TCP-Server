package znet

import (
	"TCP-Server/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	// 当前连接的套接字
	Conn *net.TCPConn
	// 当前连接的ID
	ConnID uint32
	// 当前连接是否关闭
	isClosed bool

	// 告知当前连接已经退出的channel, 连接如果要退出的话,通过管道告知一下
	ExitChan chan bool
	// 该连接处理的方法Router
	Router ziface.Irouter
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.Irouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Groutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("c.Conn.Read err = ", err)
		// 	continue
		// }

		// 创建一个拆包解包对象
		dp := NewDataPack()

		/*
			ReadFull从r精确地读取len(buf)字节数据填充进buf。
			函数返回写入的字节数和错误（如果没有读取足够的字节）。
			只有没有读取到字节时才可能返回EOF;
			如果读取了有但不够的字节时遇到了EOF, 函数会返回ErrUnexpectedEOF。
			只有返回值err为nil时，返回值n才会等于len(buf)。
		*/

		// 读取客户端的Msg Head 8字节长度
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read headData err = ", err)
			break
		}

		// 拆包, 得到msgId和msgDataLen, 放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error = ", err)
			break
		}

		// 继续根据msgDataLen读取Data, 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msgdata err = ", err)
				break
			}
		}
		msg.SetMsgData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 执行注册的路由方法
		go func(request ziface.Irequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// 启动连接, 开始当前连接的工作
func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID", c.ConnID)
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

// 发送数据, 将数据先封包, 然后发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection is closed when send msg")
	}

	// 下面将data进行封包
	// 定义一个用于封包拆包的对象
	dp := NewDataPack()

	// 封包
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id = ", msgId, "error")
		return errors.New("conn write msg error")
	}
	return nil
}
