package ziface

/*
	封包、拆包模块
	直接面向TCP连接中的数据流, 用于处理TCP粘包问题
*/

type IdataPack interface {
	// 获取数据包头部长度的方法
	GetHeadLen() uint32

	// 封包方法
	Pack(Imessage) ([]byte, error)

	// 拆包方法
	UnPack([]byte) (Imessage, error)
}
