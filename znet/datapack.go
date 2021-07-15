package znet

import (
	"TCP-Server/utils"
	"TCP-Server/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

// 封包、拆包的具体模块
type DataPack struct{}

// 封包、拆包实例的一个具体化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取数据包头部长度的方法
func (dp *DataPack) GetHeadLen() uint32 {
	// Id uint32(4字节) + DataLen uint32(4字节)
	return 8
}

// 封包方法
func (dp *DataPack) Pack(msg ziface.Imessage) ([]byte, error) {
	/*
		NewBuffer使用buf作为初始内容创建并初始化一个Buffer。
		本函数用于创建一个用于读取已存在数据的buffer;
		也用于指定用于写入的内部缓冲的大小,
		此时,buf应为一个具有指定容量但长度为0的切片。
		buf会被作为返回值的底层缓冲切片。
	*/

	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将DataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 将DataId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将数据写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法, 读取包的Head信息
func (dp *DataPack) UnPack(binaryData []byte) (ziface.Imessage, error) {
	// 创建并初始化一个二进制数据缓冲
	dataBuff := bytes.NewBuffer(binaryData)

	// 定义一个Message结构体, 用于存放解压后的数据, 只解压头部信息
	msg := &Message{}

	// 读取头部长度到msg
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读取头部ID到msg
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超过了我们配置文件中允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg dataLen recv")
	}

	return msg, nil
}
