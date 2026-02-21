package impl

import (
	"cmdb/apps/secret"

	"github.com/infraboard/mcube/v2/ioc"
	iocmongo "github.com/infraboard/mcube/v2/ioc/config/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

// 类型约束，强制要求SecretServiceImpl必须实现secret.Service的所有方法
// 使用nil初始化零开销，不构造真实对象
var _ secret.Service = (*SecretServiceImpl)(nil)

func init() {
	ioc.Controller().Registry(&SecretServiceImpl{})
}

type SecretServiceImpl struct {
	ioc.ObjectImpl
	coll *mongo.Collection
}

func (s *SecretServiceImpl) Name() string {
	return secret.AppName
}

func (s *SecretServiceImpl) Init() error {
	s.coll = iocmongo.DB().Collection("secrets")
	return nil
}
