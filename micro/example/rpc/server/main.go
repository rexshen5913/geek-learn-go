package main

import (
	"github.com/rexshen5913/geek-learn-go/geektime-go /micro/rpc"
	"github.com/rexshen5913/geek-learn-go/geektime-go /micro/rpc/serialize/json"
	"github.com/rexshen5913/geek-learn-go/geektime-go /micro/rpc/serialize/proto"
)

func main() {
	svr := rpc.NewServer()
	svr.RegisterService(&UserService{})
	svr.RegisterService(&UserServiceProto{})
	svr.RegisterSerializer(json.Serializer{})
	svr.RegisterSerializer(proto.Serializer{})
	if err := svr.Start(":8081"); err != nil {
		panic(err)
	}
}
