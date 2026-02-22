package api

import (
	"cmdb/apps/resource"
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/http/binding"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
)

// 将内容注册到IoC池子
func init() {
	ioc.Api().Registry(&ResourceApiHandler{})
}

type ResourceApiHandler struct {
	ioc.ObjectImpl
}

// 模块名字，和impl区分
func (s *ResourceApiHandler) Name() string {
	return resource.AppName + "_api"
}

// 将grpcserver进行注册通过ioc容器
func (s *ResourceApiHandler) Init() error {
	// 获取gorestful的webservice对象实例
	ws := gorestful.ObjectRouter(s)

	// mcube框架需要注入的tags标签
	tags := []string{"resource资源管理"}
	ws.Route(ws.GET("").To(s.Search).Doc("资源管理").
		// swagger api文档注解
		Param(ws.PathParameter("page_size", "分页大小").DataType("integer")).
		Param(ws.PathParameter("page_number", "页面数量").DataType("integer")).
		Param(ws.PathParameter("keywords", "关键字信息").DataType("string")).
		// tag元数据注入
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// swagger api文档生成所需内容
		Writes(resource.ResourceSet{}).
		Returns(200, "Normal", resource.ResourceSet{}).
		Returns(404, "Error", ""))
	return nil
}

func (s *ResourceApiHandler) Search(req *restful.Request, resp *restful.Response) {
	sr := resource.NewSearchRequestSet()
	// 调用mcube封装的gin bind逻辑获取前端参数
	if err := binding.Query.Bind(req.Request, sr); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "获取前端参数错误")
		return
	}
	res, err := resource.GetService().Search(req.Request.Context(), sr)
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "无法查询到内容，检查前后端逻辑")
		return
	}
	resp.WriteEntity(res)
	return
}
