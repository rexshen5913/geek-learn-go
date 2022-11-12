package demo

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type Server struct {
	services map[string]reflectionStub
}

func NewServer() *Server {
	return &Server{
		services: map[string]reflectionStub{},
	}
}

func (s *Server) Register(service Service) error {
	s.services[service.Name()] = reflectionStub{
		value: reflect.ValueOf(service),
	}
	return nil
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			// 考虑输出日志，然后返回
			// return
			continue
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				// 这里考虑输出日志
				conn.Close()
				return
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		// 读请求
		// 执行
		// 写回响应

		reqMsg, err:= ReadMsg(conn)
		if err != nil {
			return err
		}
		req := &Request{}
		err = json.Unmarshal(reqMsg, req)
		if err != nil {
			return err
		}
		// 可以考虑找到本地的服务，然后发起调用
		service, ok := s.services[req.ServiceName]
		if !ok {
			// 返回客户端一个错误信息
			return errors.New("找不到服务")
		}
		// 把参数传进去
		// context 建起来
		ctx := context.Background()
		data, err := service.invoke(ctx, req.MethodName, reqMsg)
		if err != nil {
			return err
		}
		data = EncodeMsg(data)
		_, err = conn.Write(data)
		if err != nil {
			return err
		}
	}
}

type reflectionStub struct {
	value reflect.Value
}


func (s *reflectionStub) invoke(ctx context.Context, methodName string, data []byte) ([]byte, error) {
	method := s.value.MethodByName(methodName)
	inType := method.Type().In(1)
	in := reflect.New(inType.Elem())
	err := json.Unmarshal(data, in.Interface())
	if err != nil {
		return nil, err
	}
	res := method.Call([]reflect.Value{reflect.ValueOf(ctx), in})
	if len(res) > 1 && !res[1].IsZero() {
		return nil, res[1].Interface().(error)
	}
	return json.Marshal(res[0].Interface())
}