package main

import (
	"cmdb/apps/resource"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:18080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := resource.NewRpcClient(conn)
	resp, err := client.Search(context.Background(), resource.NewSearchRequestSet())
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
