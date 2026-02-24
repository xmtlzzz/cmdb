package api

import (
	"cmdb/apps/resource"
	"cmdb/apps/secret"
	"log"
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/gorilla/websocket"
	"github.com/infraboard/mcube/v2/http/binding"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
)

func init() {
	ioc.Api().Registry(&SecretApiHandler{})
}

type SecretApiHandler struct {
	ioc.ObjectImpl
}

func (s *SecretApiHandler) Name() string {
	return secret.AppName + "_api"
}

// 将grpcserver进行注册通过ioc容器
func (s *SecretApiHandler) Init() error {
	// 获取gorestful的webservice对象实例
	ws := gorestful.ObjectRouter(s)
	// mcube框架需要注入的tags标签
	tags := []string{"secret凭证管理"}

	ws.Route(ws.POST("/").To(s.CreateSecret).Doc("凭证管理").
		// swagger api文档注解
		Param(ws.PathParameter("vendor", "云商分类").DataType("int64")).
		Param(ws.PathParameter("api_key", "API Key").DataType("string")).
		Param(ws.PathParameter("api_secret", "API密钥").DataType("string")).
		// tag元数据注入
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// swagger api文档生成所需内容
		Writes(secret.CreateSecretRequest{}).
		Returns(200, "Normal", secret.CreateSecretRequest{}).
		Returns(404, "Error", ""))

	ws.Route(ws.GET("/").To(s.QuerySecret).Doc("凭证管理").
		// swagger api文档注解
		Param(ws.PathParameter("page_size", "页面数量").DataType("int64")).
		Param(ws.PathParameter("page_number", "当前页面").DataType("int64")).
		// tag元数据注入
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// swagger api文档生成所需内容
		Writes(SecretSet{}).
		Returns(200, "Normal", SecretSet{}).
		Returns(404, "Error", ""))

	ws.Route(ws.GET("/{id}").To(s.DescribeSecret).Doc("凭证详情").
		// swagger api文档注解
		Param(ws.PathParameter("id", "云商凭证ID").DataType("string")).
		// tag元数据注入
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// swagger api文档生成所需内容
		Writes(secret.Secret{}).
		Returns(200, "Normal", secret.Secret{}).
		Returns(404, "Error", ""))

	ws.Route(ws.GET("/{id}/sync").To(s.SyncResource).Doc("凭证同步更新").
		// swagger api文档注解
		Param(ws.PathParameter("id", "云商凭证ID").DataType("string")).
		// tag元数据注入
		Metadata(restfulspec.KeyOpenAPITags, tags).
		// swagger api文档生成所需内容
		Writes([]resource.Resource{}).
		Returns(200, "Normal", []resource.Resource{}).
		Returns(404, "Error", ""))
	return nil
}

func (s *SecretApiHandler) CreateSecret(req *restful.Request, resp *restful.Response) {
	sr := secret.NewCreateSecretRequest()
	// 调用mcube封装的gin bind逻辑获取前端参数
	if err := binding.Query.Bind(req.Request, sr); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "获取前端参数错误")
		return
	}
	res, err := secret.GetService().CreateSecret(req.Request.Context(), sr)
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "无法查询到内容，检查前后端逻辑")
		return
	}
	resp.WriteEntity(res)
	return
}

type SecretSet struct {
	Total int64
	Items []*secret.Secret
}

func (s *SecretApiHandler) QuerySecret(req *restful.Request, resp *restful.Response) {
	qr := secret.NewQuerySecretRequest()
	// 调用mcube封装的gin bind逻辑获取前端参数
	if err := binding.Query.Bind(req.Request, qr); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "获取前端参数错误")
		return
	}
	res, err := secret.GetService().QuerySecret(req.Request.Context(), qr)
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "无法查询到内容，检查前后端逻辑")
		return
	}
	resp.WriteEntity(res)
	return
}

func (s *SecretApiHandler) DescribeSecret(req *restful.Request, resp *restful.Response) {
	// go-restful路劲参数获取传入的云商凭证ID
	qr := secret.NewDescribeSecretRequest(req.PathParameter("id"))
	res, err := secret.GetService().DescribeSecret(req.Request.Context(), qr)
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, "无法查询到内容，检查前后端逻辑")
		return
	}
	resp.WriteEntity(res)
	return
}

// websocket api改造
func (s *SecretApiHandler) SyncResource(req *restful.Request, resp *restful.Response) {
	failedSet := []*resource.Resource{}

	// go-restful路劲参数获取传入的云商凭证ID
	syncId := secret.NewSyncResourceRequest(req.PathParameter("id"))
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(resp, req.Request, nil)
	if err != nil {
		log.Println("websocket upgrade failed:", err)
		return
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			// 判断是否是客户端正常关闭连接
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("websocket closed by client")
			} else {
				log.Println("websocket read error:", err)
			}
			break
		}
		// 如果收到的消息为空则跳过本次循环
		if len(message) == 0 {
			log.Printf("recv empty message, type: %v, skipping", mt)
			continue
		}
		log.Printf("recv message type: %v  value: %s", mt, message)
		err = secret.GetService().SyncResource(req.Request.Context(), syncId, func(in secret.ResourceResponse) {
			// 每同步一条资源就实时推送给客户端
			if in.Success {
				conn.WriteJSON(in.Resource)
			} else {
				failedSet = append(failedSet, in.Resource)
			}
		})
		if err != nil {
			conn.WriteJSON(map[string]string{"error": "资源同步失败，请检查secret sync逻辑"})
			return
		}
		// 同步完成后发送结束标识
		conn.WriteJSON("complete")
		if len(failedSet) != 0 {
			conn.WriteJSON(failedSet)
		}
	}

}
