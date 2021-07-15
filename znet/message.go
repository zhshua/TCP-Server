package znet

type Message struct {
	DataLen uint32 // 消息的长度
	Id      uint32 // 消息的ID
	MsgData []byte // 消息的内容
}

// 创建一个消息类型
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		MsgData: data,
	}
}

// 获取消息Id
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.MsgData
}

// 设置消息Id
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// 设置消息长度
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}

// 设置消息内容
func (m *Message) SetMsgData(data []byte) {
	m.MsgData = data
}
