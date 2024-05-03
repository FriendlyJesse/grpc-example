package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	pb "example.com/m/v2/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (res *pb.HelloResponse, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("未传输token!")
	}
	mdMap := md.Copy()

	// 将 map 转换为 JSON
	jsonBytes, err := json.Marshal(mdMap)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
	var mdObj struct {
		Appid  []string `json:"appid"`
		Appkey []string `json:"appkey"`
	}
	json.Unmarshal(jsonBytes, &mdObj)

	if mdObj.Appid[0] != "Jesse" || mdObj.Appkey[0] != "123456" {
		return nil, errors.New("token 不正确!")
	}

	res = &pb.HelloResponse{
		ResponseMsg: "hello, " + req.RequestName,
	}
	return res, err
}

func main() {
	creds, _ := credentials.NewServerTLSFromFile("/home/jesse/project/grpc-study/key/test.pem", "/home/jesse/project/grpc-study/key/test.key")
	// 开启端口
	listen, _ := net.Listen("tcp", ":9090")
	// 创建 gRPC 服务
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// 在 gRPC 服务端中注册我们自己编写的服务
	pb.RegisterSayHelloServer(grpcServer, &server{})
	// 启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		log.Println("启动服务失败！err: ", err)
	}
}
