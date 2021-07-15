package ziface

// Irequest接口实际上是把客户端的 连接信息 和 请求的数据 包装到了一个Request中

type Irequest interface {
	// 得到当前连接
	GetConnection() Iconnection

	// 得到请求的消息数据
	GetData() []byte

	// 得到请求的消息的Id
	GetMsgId() uint32
}
