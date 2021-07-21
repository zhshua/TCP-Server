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
	定义世界聊天路由
*/
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.Irequest) {
	// 解析客户端传过来的protobuf格式的协议
	proto_msg := &pb.Talk{}
	if err := proto.Unmarshal(request.GetData(), proto_msg); err != nil {
		fmt.Println("Handle_proto.Unmarshal err", err)
		return
	}

	// 得到发送当前消息的玩家id
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("Handle_request.GetConnection().GetProperty err", err)
		return
	}

	// 通过玩家id得到玩家对象
	player := core.WorldMgrObj.GetOnlinePlayerByPid(pid.(int32))

	// 广播给其他玩家
	player.Talk(proto_msg.Content)
}
