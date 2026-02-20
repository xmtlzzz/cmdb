package impl

import (
	"cmdb/apps/resource"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	iocmongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	// 注册ioc对象
	ioc.Controller().Registry(&ResourceServiceImpl{})
}

type ResourceServiceImpl struct {
	// 继承IoC对象的实现
	ioc.ObjectImpl
	resource.UnimplementedRpcServer
	// 嵌入结构体解决初始化顺序导致mongodb对象空指针
	coll *mongo.Collection
}

// 模块名字
func (s *ResourceServiceImpl) Name() string {
	return resource.AppName
}

// 将grpcserver进行注册通过ioc容器
func (s *ResourceServiceImpl) Init() error {
	// IoC 配置加载完成后才初始化 MongoDB collection
	s.coll = iocmongo.DB().Collection("resources")
	// RpcServer就是PB生成的接口
	resource.RegisterRpcServer(grpc.Get().Server(), s)
	return nil
}
