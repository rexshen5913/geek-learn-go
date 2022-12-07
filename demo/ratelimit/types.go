package ratelimit

import "context"

type Limiter interface {
	Acquire(ctx context.Context, req interface{}) (interface{}, error)
	Release(resp interface{}, err error)
}

type Guardian interface {
	Allow(ctx context.Context, req interface{}) (cb func(), err error)
	AllowV1(ctx context.Context, req interface{}) (cb func(), resp interface{},  err error)
	OnRejection(ctx context.Context, req interface{}) (interface{}, error)
}

// func Limit() {
// 	var g Guardian
// 	cb, err := g.Allow(xx)
// 	if  err != nil {
// 		return g.OnRejection(ctx, req)
// 	}
// 	cb, resp, err := !g.Allow(xx)
// 	if err != nil {
// 		return resp, err
// 	}
//
// 	// 执行增长的业务逻辑
// }