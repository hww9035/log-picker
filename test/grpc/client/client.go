package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log-picker/test/grpc/pb"
)

func main() {
	conn, err := grpc.Dial(":8090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := pb.NewHwwServiceClient(conn)
	resp, _ := client.Hello(context.Background(), &pb.Request{
		Id:   123,
		Name: "hww",
	})
	fmt.Println(resp)
}
