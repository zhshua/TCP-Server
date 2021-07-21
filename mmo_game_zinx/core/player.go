package core

import (
	"TCP-Server/mmo_game_zinx/pb"
	"TCP-Server/ziface"
	"fmt"
	"math/rand"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	// 定义玩家id
	Pid int32
	// 与客户端进行通信的连接
	Conn ziface.Iconnection
	// 玩家的XYZV坐标
	X float32
	Y float32
	Z float32
	V float32
}

// 用于生成玩家ID的计数器
var PidGen int32 = 1
var IdLock sync.Mutex

// 初始化创建一个Player
func NewPlayer(conn ziface.Iconnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), // 随机在160坐标点基于X轴偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), // 随机在140坐标点基于Y轴偏移
		V:    0,
	}
}

/*
	提供一个给客户端发送消息的方法
	主要是将pb的protobuf数据序列化之后, 再调用zinx的SendMsg
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	// 将proto.Message结构体序列化, 转化为二进制数据
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("proto.Marshal err", err)
		return
	}

	// 调用zinx的SendMsg方法, 将序列化好的msg发送给客户端
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("p.Conn.SendMsg err", err)
		return
	}
}

// 告知客户端玩家id, 同步已经生成的玩家id给客户端
func (p *Player) SyncPid() {
	// 组建 MsgId:1 的proto数据
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	// 将消息发送给客户端
	p.SendMsg(1, proto_msg)
}

// 广播玩家自己的出生点
func (p *Player) BroadCastStartPosition() {
	// 组建 MsgId:200 的proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 将消息发送给客户端
	p.SendMsg(200, proto_msg)
}

func (p *Player) Talk(content string) {
	// 组建MsgId:200的消息格式
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	// 得到当前世界在线玩家
	allOnlinePlayers := WorldMgrObj.GetAllOnlinePlayers()

	// 向所有在线玩家(包括自己)广播发送MsgId:200的消息
	for _, player := range allOnlinePlayers {
		player.SendMsg(200, proto_msg)
	}
}

// 同步自己的位置信息给周围九宫格内的玩家
func (p *Player) SyncSurrounding() {
	// 获取周围九宫格内的所有玩家id
	pids := WorldMgrObj.AoiMgr.GetPidsbyPos(p.X, p.Z)
	surPlayers := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		surPlayers = append(surPlayers, WorldMgrObj.GetOnlinePlayerByPid(int32(pid)))
	}

	// 将当前玩家的位置信息通过MsgID:200广播发给周围的玩家
	// 1. 组装MsgID:200消息
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 2.周围玩家给各自的客户端发送玩家上线的MsgID:200消息(让其他玩家看到自己)
	for _, player := range surPlayers {
		player.SendMsg(200, proto_msg)
	}

	// 将周围玩家的信息发送给当前玩家(让自己看到其他玩家)
	// 组装MsgID:202消息格式的数据
	players_proto_msg := make([]*pb.Player, 0, len(surPlayers))
	for _, player := range surPlayers {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}
	// 同步玩家信息的消息格式202
	SyncPlayers_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}
	// 发送202格式的消息协议
	p.SendMsg(202, SyncPlayers_proto_msg)
}
