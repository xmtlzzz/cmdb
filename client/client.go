package main

import (
	"cmdb/apps/resource"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GrpcClientTest(address string) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

var addr = flag.String("addr", "localhost:8080", "http service address")

func WebSocketClientTest() {
	url := url.URL{Scheme: "ws", Host: *addr, Path: "/api/cmdb/1.0.0/secret_api/88b46c44-bfb8-3ec1-b148-fa400a18a605/sync"}
	log.Printf("connecting to %s", url.String())

	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	c.WriteMessage(websocket.TextMessage, []byte("start sync"))

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		// 服务端使用WriteJSON发送"complete"，需要先反序列化JSON
		var msg string
		if err := json.Unmarshal(message, &msg); err == nil && msg == "complete" {
			fmt.Println("sync complete")
			break
		}
		fmt.Println(string(message))
	}
}

func main() {
	//GrpcClientTest("127.0.0.1:18080")
	WebSocketClientTest()
}
