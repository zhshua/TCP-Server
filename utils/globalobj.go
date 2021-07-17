package utils

import (
	"TCP-Server/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
	存储一切有关Zinx框架的全局参数, 供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置
*/
type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的ip地址
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	/*
		Zinx
	*/
	Version          string // 当前Zinx的版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作池workerPool中goroutine的数量
	MaxWorkerTaskLen uint32 // 每个Worker对应的消息队列中任务的最大数量
}

// 定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

/*
	从已有的zinx.json文件中加载相关配置
*/
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}
}

/*
	提供一个init方法, 初始化当前的GlobalObject对象
*/
func init() {
	// 如果配置文件没有配置相关项, 则提供一个默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "Zinx-0.9",
		Host:             "0.0.0.0",
		TcpPort:          9190,
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 从配置文件里加载相关项
	GlobalObject.Reload()
}
