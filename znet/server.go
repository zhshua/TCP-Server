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
	// 当前的Server添加一个Router, server注册的连接对应的处理业务
	Router ziface.Irouter
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

			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			go dealConn.Start()

		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
}

// 路由功能：给当前的服务注册一个路由方法, 供客户端的连接处理使用
func (s *Server) AddRouter(router ziface.Irouter) {
	s.Router = router
	fmt.Println("add router succ!")
}

// 创建一个服务器句柄
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
