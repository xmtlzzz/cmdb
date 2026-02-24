package resource

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/validator"
)

const (
	AppName = "resource"
)

// 对外用于获取resource.Service接口的构造方法
func GetService() Service {
	// 断言类型为resource.Service
	return ioc.Controller().Get(AppName).(Service)
}

// 扩展方法，deleteresource为内部使用
type Service interface {
	RpcServer
	DeleteResource(context.Context, *DeleteResourceRequest) error
}

func NewDeleteResourceRequest() *DeleteResourceRequest {
	return &DeleteResourceRequest{}
}

type DeleteResourceRequest struct {
	ResourceId []string
}

// 校验字段
func (s *Resource) Validate() error {
	if err := validator.Validate(s); err != nil {
		return err
	}
	return nil
}

// 构造方法初始化ResourceSet
func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		Items: []*Resource{},
	}
}

// 排除前n页的总共n项内容
func (s *SearchRequestSet) SetSkip() int64 {
	return (s.PageNumber - 1) * s.PageSize
}

func NewSearchRequestSet() *SearchRequestSet {
	return &SearchRequestSet{
		PageSize:   10,
		PageNumber: 1,
		Tags:       map[string]string{},
	}
}

func NewResource() *Resource {
	return &Resource{
		Meta:   &Meta{},
		Spec:   &Spec{},
		Status: &Status{},
	}
}
