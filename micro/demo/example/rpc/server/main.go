package main

import (
	rpc "github.com/rexshen5913/geek-learn-go/geektime-go /micro/demo"
	"github.com/rexshen5913/geek-learn-go/geektime-go /micro/demo/serialize/json"
	"github.com/rexshen5913/geek-learn-go/geektime-go /micro/demo/serialize/proto"
)

func main() {
	svr := rpc.NewServer()
	svr.MustRegister(&UserService{})
	svr.MustRegister(&UserServiceProto{})
	svr.RegisterSerializer(json.Serializer{})
	svr.RegisterSerializer(proto.Serializer{})
	if err := svr.Start(":8081"); err != nil {
		panic(err)
	}
}
