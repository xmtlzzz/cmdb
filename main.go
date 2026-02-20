package main

import (
	"cmdb/test"
	"context"
	"log"

	_ "cmdb/apps"

	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/restful"
	"github.com/infraboard/mcube/v2/ioc/server"
)

func init() {
	test.SetUp()
}

func main() {
	if err := server.Run(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
