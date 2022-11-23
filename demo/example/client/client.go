package main

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/demo"
	"gitee.com/geektime-geekbang/geektime-go/demo/example/proto/gen"
	"google.golang.org/grpc"
	"log"
)

func main() {
	rsBuilder := demo.NewResolverBuilder()
	cc, err := grpc.Dial("registry:///user-service",
		grpc.WithInsecure(),
		grpc.WithResolvers(rsBuilder))
	if err != nil {
		panic(err)
	}
	usClient := gen.NewUserServiceClient(cc)
	resp, err := usClient.GetById(context.Background(), &gen.GetByIdReq{
		Id: 12,
	})
	if err != nil {
		panic(err)
	}
	log.Println(resp)
}
