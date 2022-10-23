package demo

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/demo/internal/errs"
	"sync"
	"time"
)



type LocalCache struct {
	m sync.Map
}

func NewLocalCache() *LocalCache {
	res := &LocalCache{}
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <- ticker.C:
				res.m.Range(func(key, value any) bool {
					// 如果过期了
					res.m.Delete(key)
					return true
				})

			case :
				// 监听关闭的信号

			}
		}
	}()
	return res
}

func (l *LocalCache) Get(ctx context.Context, key string) (any, error) {
	val, ok := l.m.Load(key)
	if !ok {
		return nil, errs.NewErrKeyNotFound(key)
	}
	return val, nil
}

func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	// time.AfterFunc(expiration, func() {
	// 	l.Delete(context.Background(), key)
	// })
	l.m.Store(key, val)
	return nil
}

func (l *LocalCache) Delete(ctx context.Context, key string) error {
	l.m.Delete(key)
	return nil
}

