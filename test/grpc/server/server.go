package main

import (
	"context"
	"google.golang.org/grpc"
	"log-picker/test/grpc/pb"
	"net"
)

type HwwServer struct {
	pb.UnimplementedHwwServiceServer
}

func (s *HwwServer) Hello(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	return &pb.Response{
		Msg:  "server: " + request.Name,
		Code: pb.Stat(request.Id),
	}, nil
}

func (s *HwwServer) mustEmbedUnimplementedHwwServiceServer() {
	//TODO implement me
	panic("implement me")
}

func main() {
	// 实例化grpc服务
	ser := grpc.NewServer()

	// 注册服务
	pb.RegisterHwwServiceServer(ser, new(HwwServer))

	// 监听并运行服务
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic("rpc listen fail.")
	}
	err = ser.Serve(lis)
	if err != nil {
		panic("rpc listen fail.")
	}
}
