package main

import (
	"cmdb/test"
	"context"
	"log"

	_ "cmdb/apps"

	"github.com/infraboard/mcube/v2/ioc"
	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/restful"
	"github.com/infraboard/mcube/v2/ioc/server"
)

func init() {
	test.SetUp()
	// 去使能debug日志
	ioc.SetDebug(false)
}

func main() {
	if err := server.Run(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
