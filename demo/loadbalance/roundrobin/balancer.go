package roundrobin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync/atomic"
)

type Picker struct {
	conns []balancer.SubConn
	cnt uint64
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	// 这个不是线程安全操作
	// p.cnt = p.cnt + 1


	cnt := atomic.AddUint64(&p.cnt, 1)
	index := cnt % uint64(len(p.conns))
	return balancer.PickResult{
		SubConn: p.conns[index],
		Done: func(info balancer.DoneInfo) {
			// 这个地方就是很神奇的地方
			// 效果就是根据调用结果来调整你的负载均衡策略

			// 假如说出错了
			// if info.Err != nil {
				// 尝试把这个 subConn 置为不可用，或者临时移除出去
			// }
		},
	}, nil
}

type PickerBuilder struct {

}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for conn := range info.ReadySCs {
		conns = append(conns, conn)
	}
	return &Picker{
		conns: conns,
	}
}

func (p *PickerBuilder) Name() string {
	return "ROUND_ROBIN"
}



