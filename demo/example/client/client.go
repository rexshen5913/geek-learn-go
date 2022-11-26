package main

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/demo"
	"gitee.com/geektime-geekbang/geektime-go/demo/example/proto/gen"
	"gitee.com/geektime-geekbang/geektime-go/demo/registry"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var r registry.Registry
	rsBuilder := demo.NewResolverBuilder(r, time.Second)
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
