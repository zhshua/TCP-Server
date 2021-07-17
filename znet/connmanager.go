package znet

import (
	"TCP-Server/ziface"
	"errors"
	"fmt"
	"sync"
)

/*
	连接管理模块的实现
*/
type ConnManager struct {
	// 存储ConnID与对应连接的集合
	connetcions map[uint32]ziface.Iconnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

// 创建当前连接管理模块
func NewConnManager() *ConnManager {
	return &ConnManager{
		connetcions: make(map[uint32]ziface.Iconnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.Iconnection) {
	// 添加的时候要上写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到connManager中
	connMgr.connetcions[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), " add to ConnManager succ : conn nums = ", connMgr.Len())

}

// 删除连接
func (connMgr *ConnManager) Remove(conn ziface.Iconnection) {
	// 删除时候上写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connetcions, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), "remove to ConnManager succ : conn nums = ", connMgr.Len())

}

// 根据ConnID获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.Iconnection, error) {
	// 获取连接加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connetcions[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("conn not found")
}

// 得到连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connetcions)
}

// 清除所有连接
func (connMgr *ConnManager) ClearConn() {
	// 删除连接加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除所有连接, 并停止连接的工作
	for connID, conn := range connMgr.connetcions {
		conn.Stop()
		delete(connMgr.connetcions, connID)
	}
	fmt.Println("Clear All conns succ! conn nums = ", connMgr.Len())
}
