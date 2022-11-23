package demo

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/demo/registry"
	"google.golang.org/grpc/resolver"
)

type grpcResolverBuilder struct {
	r registry.Registry
}

func NewResolverBuilder() resolver.Builder {
	return &grpcResolverBuilder{

	}
}

func (g *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	eventsCh, err := g.r.Subscribe(target.Endpoint)
	if err != nil {
		return nil, err
	}
	res := &grpcResolver{
		target: target,
		cc: cc,
		r: g.r,
	}

	go func() {
		for {
			select {
			case <-eventsCh:
				// 立刻更新可用节点列表
				res.resolve()
			}
		}
	}()

	return res, nil
}

func (g *grpcResolverBuilder) Scheme() string {
	return "registry"
}

type grpcResolver struct {
	target resolver.Target
	cc resolver.ClientConn
	r registry.Registry
}

// ResolveNow 立刻解析——立刻执行服务发现——立刻去问一下注册中心
func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.resolve()
}

func (g *grpcResolver) resolve() {
	r := g.r
	// 这个就是可用服务实例（节点）列表
	instances, err := r.ListService(context.Background(), g.target.Endpoint)
	if err != nil {
		g.cc.ReportError(err)
		return
	}
	// 我是不是要把 instances 转成 Address
	address := make([]resolver.Address, 0, len(instances))
	for _, ins := range instances {
		address = append(address, resolver.Address{
			// 定位信息，ip+端口
			Addr: ins.Address,
		})
	}
	err = g.cc.UpdateState(resolver.State{
		Addresses: address,
	})
	if err != nil {
		g.cc.ReportError(err)
	}
}

func (g *grpcResolver) Close() {
	// TODO implement me
	panic("implement me")
}

