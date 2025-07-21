package rpc

import (
	"context"
	"github.com/rexshen5913/geek-learn-go/geektime-go /micro/rpc/message"
)

type Proxy interface {
	Invoke(ctx context.Context, req *message.Request) (*message.Response, error)
}

type Service interface {
	ServiceName() string
}
