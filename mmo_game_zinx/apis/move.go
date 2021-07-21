package apis

import (
	"TCP-Server/mmo_game_zinx/core"
	"TCP-Server/mmo_game_zinx/pb"
	"TCP-Server/ziface"
	"TCP-Server/znet"
	"fmt"

	"google.golang.org/protobuf/proto"
)

/*
	定义玩家移动路由
*/
type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.Irequest) {
	// 解析客户端发来的data
	proto_msg := &pb.Position{}
	if err := proto.Unmarshal(request.GetData(), proto_msg); err != nil {
		fmt.Println("proto.Unmarshal err", err)
		return
	}
	// 得到发来位置信息的那个玩家id
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("request.GetConnection().GetProperty() err", err)
		return
	}
	fmt.Printf("Player pid = %d, move(%f,%f,%f,%f)\n", pid.(int32), proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	// 通过玩家id获得玩家对象
	player := core.WorldMgrObj.GetOnlinePlayerByPid(pid.(int32))

	// 广播自己当前的位置信息给其他玩家
	player.UpdataPos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
