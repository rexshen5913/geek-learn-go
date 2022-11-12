package demo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/silenceper/pool"
	"net"
	"time"
)

type Client struct {
	connPool pool.Pool
}

func NewClient(addr string) (*Client, error){
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap: 10,
		MaxCap: 100,
		MaxIdle: 50,
		Factory: func() (interface{}, error) {
			return net.Dial("tcp", addr)
		},
		IdleTimeout: time.Minute,
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		connPool: p,
	}, nil
}

func (c *Client) Invoke(ctx context.Context, req *Request) (*Response, error) {
	// 发送请求过去
	data, err:= json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// 拿一个连接
	obj, err := c.connPool.Get()
	// 这个 error 是框架 error，而不是用户返回的 error
	if err != nil {
		return nil, err
	}
	conn := obj.(net.Conn)
	// 发请求

	// 0001 0002 -- 这一坨是描述你数据有多长
	// 0000  --- 这一坨是数据
	data = EncodeMsg(data)
	i, err := conn.Write(data)
	if err != nil {
		return  nil, err
	}
	if i != len(data) {
		return nil, errors.New("micro: 未写入全部数据")
	}
	// 读响应
	// 我怎么知道该读多长数据？相应地，服务端读请求，该读多长？
	// 先读长度字段

	// 读取全部的响应
	// 装响应的 bytes
	respMsg, err := ReadMsg(conn)
	if err != nil {
		return nil, err
	}
	return &Response{
		Data: respMsg,
	}, nil
}

// 客户端
// 代码演示第一部分
// 1. 首先反射拿到 Request，核心是服务名字，方法名字和参数
// 2. 将 Request 进行编码，要注意序列化并且加上长度字段
// 3. 使用连接池，或者一个连接，把请求发过去

// 代码演示第四部分
// 4. 从连接里面读取响应，解析成结构体

// 服务端
// 代码演示第二部分
// 1. 启动一个服务器，监听一个端口
// 2. 读取长度字段，再根据长度，读完整个消息
// 3. 解析成 Request
// 4. 查找服务，会对应的方法
// 5. 构造方法对应的输入
// 6. 反射执行调用

// 代码演示第三部分
// 7. 编码响应
// 8. 写回响应

