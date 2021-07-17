package ziface

/*
	消息管理抽象层
*/

type IMsgHandle interface {
	// 调度/执行对应的Router消息处理方法
	DoMsgHandler(Irequest)
	// 为消息添加具体的处理逻辑
	AddRouter(uint32, Irouter)
	// 开启工作池
	StartWorkPool()
	// 将消息发送给消息任务队列处理
	SendMsgToTaskQueue(Irequest)
}
