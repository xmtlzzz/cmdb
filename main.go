package main

import (
	_ "cmdb/apps"
	"cmdb/test"

	"github.com/infraboard/mcube/v2/ioc"
	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/restful"
	"github.com/infraboard/mcube/v2/ioc/server/cmd"
)

func init() {
	test.SetUp()
	// 去使能debug日志
	ioc.SetDebug(false)
}

func main() {
	cmd.Start()
}
