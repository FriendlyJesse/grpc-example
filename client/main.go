package main

import (
	"context"
	"fmt"
	"log"

	pb "example.com/m/v2/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientTokenAuth struct {
}

func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (data map[string]string, err error) {
	data = map[string]string{
		"appid":  "Jesse",
		"appkey": "123456",
	}

	return data, err
}
func (c ClientTokenAuth) RequireTransportSecurity() bool {
	return true
}

func main() {
	creds, _ := credentials.NewClientTLSFromFile("/home/jesse/project/grpc-study/key/test.pem", "*.test.com")
	// 连接到 server，此处禁用安全加密，没有加密和验证
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(new(ClientTokenAuth)))
	if err != nil {
		log.Fatalln("did not connect:", err)
	}
	defer conn.Close()

	// 建立连接
	client := pb.NewSayHelloClient(conn)
	ctx := context.Background()
	// 执行 rpc
	res, err := client.SayHello(ctx, &pb.HelloRequest{RequestName: "Jesse"})
	if err != nil {
		log.Println("rpc err: ", err)
	}
	fmt.Println(res.GetResponseMsg())
}
