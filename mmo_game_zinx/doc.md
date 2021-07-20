- apis: 用来存放基本的用户自定义路由业务, 一个msgID对应的一个业务
- conf: 存放zinx的配置文件
- pb: msg.proto原始的protobuf协议文件
      build.sh 编译msg.proto的脚本
      msg.pb.go 编译生成的go文件(只读)
- core: 存放核心的功能
- main.go: 服务器的入口函数