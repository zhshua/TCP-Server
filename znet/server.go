package znet

import (
	"TCP-Server/utils"
	"TCP-Server/ziface"
	"fmt"
	"net"
	"time"
)

type Server struct {
	// 服务器名称
	Name string
	// 服务器的IP地址版本 tcp4 or other
	IPVersion string
	// 服务器的IP地址
	IP string
	// 监听的端口号
	Port int
	// 当前Server的消息管理模块, 用来绑定MsgID和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle
	// 该Server的连接管理器
	ConnMgr ziface.IConnManager
	// 该Server创建连接之后自动调用的Hook函数
	OnConnStart func(conn ziface.Iconnection)
	// 该Server销毁连接之前自动调用的Hook函数
	OnConnStop func(conn ziface.Iconnection)
}

// 运行服务器
func (s *Server) Serve() {
	// 调用Serve()启动服务器
	s.Start()
	// 这里之所以对Serve封装了一层Start,
	// 是因为在服务器启动后我们可能还会去做一些额外的事情

	//TODO Server.Serve()是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	// 阻塞一下, 不让Serve退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// 运行服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s, listenner at IP : %s, Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPackageSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	fmt.Printf("[START] Server listen at ip:%s, port:%d\n", s.IP, s.Port)

	// 开启一个goroutine去监听服务端的lister业务
	go func() {
		// 开启Worker工作池
		s.MsgHandle.StartWorkPool()

		// 绑定监听地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("net.ResolveTCPAddr err = ", err)
			return
		}

		// 开始监听
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("net.ListenTCP err = ", err)
			return
		}
		fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")

		// 阻塞等待接受连接
		var cid uint32 = 0
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("listenner.AcceptTCP err = ", err)
				continue
			}

			// 当前最大连接个数判断, 如果超过个数, 则关闭此连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO：给客户端响应一个超出最大连接的错误包
				fmt.Println("too many conns maxconns = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid++

			go dealConn.Start()

		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	s.ConnMgr.ClearConn()
}

// 路由功能：给当前的服务注册一个路由方法, 供客户端的连接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.Irouter) {
	s.MsgHandle.AddRouter(msgId, router)
	fmt.Println("add router succ!")
}

// 创建一个服务器句柄
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}
	return s
}

// 获取连接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册 OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.Iconnection)) {
	s.OnConnStart = hookFunc
}

// 注册 OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.Iconnection)) {
	s.OnConnStop = hookFunc
}

// 调用 OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(connection ziface.Iconnection) {
	// hookfunc注册之后才能调用
	if s.OnConnStart != nil {
		fmt.Println("-----> Call OnConnStart()...")
		s.OnConnStart(connection)
	}
}

// 调用 OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(connection ziface.Iconnection) {
	// hookfunc注册之后才能调用
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStop()...")
		s.OnConnStop(connection)
	}
}
