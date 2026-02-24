package impl_test

import (
	"cmdb/apps/resource"
	_ "cmdb/apps/resource/impl"
	"cmdb/apps/secret"
	_ "cmdb/apps/secret/impl"
	"cmdb/test"
	"context"
	"os"
	"testing"
)

var (
	ctx context.Context
	svc = secret.GetService()
)

func init() {
	test.SetUp()
}

func TestSecretServiceImpl_CreateSecret(t *testing.T) {
	ins := secret.NewCreateSecretRequest()
	ins.Name = "阿里云只读用户"
	ins.Vendor = resource.Vendor_ALIYUN
	ins.ApiKey = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	ins.ApiSecret = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	ins.Regions = []string{"SH", "BJ"}
	se, err := svc.CreateSecret(ctx, ins)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(se)
}

func TestSecretServiceImpl_QuerySecret(t *testing.T) {
	qs := secret.NewQuerySecretRequest()
	res, err := svc.QuerySecret(ctx, qs)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestSecretServiceImpl_DescribeSecret(t *testing.T) {
	ds := secret.NewDescribeSecretRequest("88b46c44-bfb8-3ec1-b148-fa400a18a605")
	res, err := svc.DescribeSecret(ctx, ds)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestSecretServiceImpl_SyncResource(t *testing.T) {
	if err := svc.SyncResource(ctx, secret.NewSyncResourceRequest("88b46c44-bfb8-3ec1-b148-fa400a18a605"), func(in secret.ResourceResponse) {
		t.Log(in)
	}); err != nil {
		t.Fatal(err)
	}
}
