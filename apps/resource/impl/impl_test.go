package impl_test

import (
	"cmdb/apps/resource"
	"cmdb/test"
	"context"
)

var (
	ctx = context.Background()
	svc = resource.GetService()
)

func init() {
	test.SetUp()
}
