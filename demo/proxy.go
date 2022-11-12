package demo

import "context"

type Proxy interface {
	Invoke(ctx context.Context, req *Request) (*Response, error)
}

type Request struct {
	ServiceName string
	MethodName string
	Data []byte
}

type Response struct {
	Data []byte
}