package secret

import (
	"cmdb/apps/resource"
	local_errors "errors"
	"fmt"

	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (s *Secret) Sync(cb SyncResourceHandleFunc) error {

	switch s.Vendor {
	case resource.Vendor_TENCENT:
		// 配置api信息
		credential := common.NewCredential(
			s.ApiKey,
			s.ApiSecret,
		)
		// 使用临时密钥示例
		// credential := common.NewTokenCredential("SecretId", "SecretKey", "Token")
		// 实例化一个client选项，可选的，没有特殊需求可以跳过
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"

		for _, region := range s.Regions {
			// 实例化要请求产品的client对象,clientProfile是可选的
			client, _ := cvm.NewClient(credential, region, cpf)

			// 实例化一个请求对象,每个接口都会对应一个request对象
			request := cvm.NewDescribeInstancesRequest()

			// 设置request的limit、offset,每个region都会分别进行请求
			SetLimit(request, 10)
			SetOffset(request, 0)

			// 循环遍历所有内容
			for {
				response, err := client.DescribeInstances(request)
				if _, ok := err.(*errors.TencentCloudSDKError); ok {
					return local_errors.New(fmt.Sprintf("An API error has returned: %s", err))
				}
				if err != nil {
					return err // 也建议改为 return err，而非 panic
				}
				// 遍历Instance并且通过自定义函数格式化内容到resource.Resource
				for _, ins := range response.Response.InstanceSet {
					cb(ResourceResponse{Resource: FormatTencentCVM(ins)})
				}
				// 正确终止条件：已取完所有数据
				*request.Offset += int64(len(response.Response.InstanceSet))
				if *request.Offset >= *response.Response.TotalCount {
					// json格式化
					pretty.ToJSON(response)
					break
				}
			}
		}
	case resource.Vendor_ALIYUN:
	}
	return nil
}

// 将cvm主机的信息格式化
func FormatTencentCVM(ins *cvm.Instance) *resource.Resource {
	res := resource.NewResource()
	// 具体的转化逻辑
	res.Meta.Id = GetValue(ins.InstanceId)
	res.Spec.Name = GetValue(ins.InstanceName)
	res.Spec.Cpu = GetValue(ins.CPU)
	res.Spec.Memory = GetValue(ins.Memory)
	res.Spec.Storage = GetValue(ins.SystemDisk.DiskSize)
	res.Status.PrivateAddress = common.StringValues(ins.PrivateIpAddresses)
	return res
}

// 如果存在值就传递，不存在就传递空值
func GetValue[T any](in *T) T {
	if in == nil {
		var v T
		return v
	}
	return *in
}

// 设置limit，分页
func SetLimit(request *cvm.DescribeInstancesRequest, in int64) {
	request.Limit = &in
}

// 设置偏移量，分页
func SetOffset(request *cvm.DescribeInstancesRequest, in int64) {
	request.Offset = &in
}
