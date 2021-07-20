package znet

import (
	"TCP-Server/utils"
	"TCP-Server/ziface"
	"fmt"
	"strconv"
)

/*
	消息处理模块的实现
*/
type MsgHandle struct {
	// 不同的消息对应不同的路由处理
	Apis map[uint32]ziface.Irouter
	// Worker负责从消息队列中取任务, 一个Worker对应一个消息队列
	TaskQueue []chan ziface.Irequest
	// 工作池中Worker的数量
	WorkerPoolSize uint32
}

// 初始化MsgHandle
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.Irouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.Irequest, utils.GlobalObject.WorkerPoolSize),
	}
}

/* 调度/执行对应的Router消息处理方法 */
func (mh *MsgHandle) DoMsgHandler(request ziface.Irequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), "not found!")
		return
	}

	// 调度对应Router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.Irouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeat api, msgId = " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("add api msgId = ", msgId, " succ ")
}

// 启动一个工作池, 开启工作池的动作只能发生一次, 一个zinx框架只能有一个Worker工作池
func (mh *MsgHandle) StartWorkPool() {
	// 根据WorkPoolSize 分别开启Worker, 每个Worker用一个goroutine承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.Irequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前的Worker, 阻塞等待消息从Channel发过来
		// 一个Worker对应一个消息任务处理队列
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(WorkerID int, taskQueue chan ziface.Irequest) {
	fmt.Println("Worker id = ", WorkerID, "is started...")

	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息连接过来, 出列的就是一个客户端的Request, 则执行当前Request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue, 由Worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.Irequest) {
	// 将消息平均分配给某个worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnId = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgId(),
		" to WorkerID = ", workerID)

	// 将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
