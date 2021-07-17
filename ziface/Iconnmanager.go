package ziface

/*
	连接管理模块
*/

type IConnManager interface {
	// 添加连接
	Add(conn Iconnection)
	// 删除连接
	Remove(conn Iconnection)
	// 根据ConnID获取连接
	Get(connID uint32) (Iconnection, error)
	// 得到连接总数
	Len() int
	// 清除所有连接
	ClearConn()
}
