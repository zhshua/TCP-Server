package ziface

// 路由抽象接口, 路由里的数据都是Irequest
type Irouter interface {
	// 在处理conn业务之前的钩子方法Hook
	PreHandle(request Irequest)

	// 处理conn业务的主方法 Hook
	Handle(request Irequest)

	// 在处理conn业务之后的钩子方法Hook
	PostHandle(request Irequest)
}
