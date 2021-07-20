## 文件分布
- apis: 用来存放基本的用户自定义路由业务, 一个msgID对应的一个业务
- conf: 存放zinx的配置文件
- pb: 
  - msg.proto原始的protobuf协议文件
  - build.sh 编译msg.proto的脚本
  - msg.pb.go 编译生成的go文件(只读)
- core: 存放核心的功能
- main.go: 服务器的入口函数

## 协议制定
### MsgID:1
```
message SyncPid{
      int32 Pid=1;
}
```
- SyncPid
  - 同步玩家本次登录的ID(用来标识玩家),玩家登录后,由Server端主动生成玩家ID发送给客户端
  - 发起者: Server
  - Pid: 玩家ID

### MsgID:2
- Talk
  - 玩家聊天时候发送的信息
  - 发起者: Client
  - Content: 聊天信息
```
message Talk{
      string Content = 1;
}
```
### MsgID:3
```
message Position{
    float X = 1;    // 平面的X坐标
    float Y = 2;    // 高度
    float Z = 3;    // 平面的Y坐标(注意这里的Y和Z和平常坐标系的设计不一样)
    float V = 4;    // 玩家的倾斜角度(0-360°)
}
```
- MovePackage
  - 移动的坐标数据
  - 发起者: Client
  - P: Position类型, 地图的坐标点
  - X,Y,Z三维坐标, V表示玩家倾斜角度, 比如上坡
### MsgID:200
```
message BroadCast{
    int32 Pid = 1;
    int32 Tp = 2;   
    // oneof表示下面三种类型只能取一个
    oneof Data{
        string Content = 3;     // 玩家聊天信息
        Position P = 4;         // 广播玩家位置
        int32 ActionData = 5;   // 玩家具体动作
    }
}
```
- BroadCast
  - 广播消息
  - 发起者: Client
  - Tp=1:世界聊天; Tp=2:玩家位置; Tp=3:动作; Tp=4:移动之后的坐标信息
  - Pid: 玩家ID
### MsgID:201
```
message SyncPid{
      int32 Pid = 1;
}
```
- SyncPid
  - 广播消息, 玩家掉线, 或者消失在视野中
  - 发起者: Server
  - Pid: 玩家id
### MsgID:202
```
message SyncPlayers{
      repeated Player ps = 1;
}
message Player{
      int32 Pid = 1;
      Position P = 2;
}
```
- SyncPlayers
  - 同步周围人的信息位置(包括自己)
  - 发起者: Server
  - ps: Player集合, 需要同步的玩家



