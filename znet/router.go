package znet

import "TCP-Server/ziface"

// 实现router时，先嵌入这个BaseRouter类，然后根据需要对这个基类的方法进行重写就好了
type BaseRouter struct{}

// 这里之所以BaseRouter的方法都为空
// 是因为有的Router不希望有PreHandle、PostHandle这两个业务
// 所以Router全部继承BaseRouter的好处就是, 不需要实现PreHandle、PostHandle业务

// 在处理conn业务之前的钩子方法Hook
func (rb *BaseRouter) PreHandle(request ziface.Irequest) {}

// 处理conn业务的主方法 Hook
func (rb *BaseRouter) Handle(request ziface.Irequest) {}

// 在处理conn业务之后的钩子方法Hook
func (rb *BaseRouter) PostHandle(request ziface.Irequest) {}
