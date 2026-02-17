package main

import (
	"context"
	"log"

	"github.com/infraboard/mcube/v2/ioc/server"
)

func main() {
	if err := server.Run(context.TODO()); err != nil {
		log.Fatal(err)
	}
}
