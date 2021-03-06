package znet

import (
	"TCP-Server/utils"
	"TCP-Server/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	// 当前Conn隶属于哪个Server
	TcpServer ziface.IServer
	// 当前连接的套接字
	Conn *net.TCPConn
	// 当前连接的ID
	ConnID uint32
	// 当前连接是否关闭
	isClosed bool

	// 告知当前连接已经退出的channel, 连接如果要退出的话,通过管道告知一下
	ExitChan chan bool

	// 无缓冲管道, 用于读goroutine和写goroutine之间的消息通信
	msgChan chan []byte

	// 消息管理模块 MsgId和对应的处理业务API的关系
	MsgHandle ziface.IMsgHandle

	// 连接属性集合
	property map[string]interface{}

	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

// 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		MsgHandle: msgHandle,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
		property:  make(map[string]interface{}),
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Groutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
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
		// 读完放在msg.Data中
		msg.SetMsgData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 判断是否已经开启工作池, 如果开启了工作池机制, 则将消息发送给Worker工作池处理
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中找到注册绑定的Conn对应的Router调用
			go c.MsgHandle.DoMsgHandler(&req)
		}

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("Reader Groutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Writer is exit, remote addr is ", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出, Writer也要退出
			return
		}
	}
}

// 启动连接, 开始当前连接的工作
func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID", c.ConnID)
	go c.StartReader()
	go c.StartWriter()

	// 按照开发者传递进来的 创建连接之后需要调用的处理业务, 执行对应Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止连接, 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果连接已经是关闭的
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 按照开发者传递进来的 关闭连接之前需要执行的业务, 调用Hook函数
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从连接管理器中去除
	c.TcpServer.GetConnMgr().Remove(c)

	close(c.ExitChan)
	close(c.msgChan)
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
	if c.isClosed {
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

	c.msgChan <- binaryMsg
	return nil
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	// 加写锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	// 加读锁
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no peoperty found")
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	// 加写锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
