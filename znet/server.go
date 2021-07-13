package znet

import (
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

		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("listenner.AcceptTCP err = ", err)
				continue
			}

			// 启动一个goroutine 分离出读写模块
			go func() {
				for {
					// 读取数据
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("conn.Read err = ", err)
						continue
					}
					fmt.Printf("recv %d byte data: '%s' from %s\n", cnt, buf[:cnt], conn.RemoteAddr())

					// 回显
					if _, err = conn.Write(buf[:cnt]); err != nil {
						fmt.Println("conn.Write err = ", err)
						continue
					}
				}
			}()

		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
}

// 创建一个服务器句柄
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      9190,
	}
	return s
}
