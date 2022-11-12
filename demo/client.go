package demo

import (
	"context"
	"encoding/json"
	"reflect"
)

// func InitClientProxyV1[T Service](service T) error {
//
// }

func InitClientProxy(service Service, p Proxy) error {
	// 你可以做校验，确保它必须是一个指向结构体的指针

	val := reflect.ValueOf(service).Elem()
	typ := reflect.TypeOf(service).Elem()
	numField := val.NumField()
	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			// 可以报错，也可以跳掉
			continue
		}
		// if fieldType.Type.Kind() != reflect.Func {
		// 	continue
		// }
		// 替换新的实现
		// 替换为一个新的方法实现
		fn := reflect.MakeFunc(fieldType.Type,
			func(args []reflect.Value) (results []reflect.Value) {
				// 实际上你在这里需要对 args 和 results 进行校验
				// 第一个返回值，真的返回值，指向 GetIdResp
				outType := fieldType.Type.Out(0)
				ctx := args[0].Interface().(context.Context)
				// if !ok {
				// 	return errors.Ne
				// }
				arg := args[1].Interface()

				bs, err := json.Marshal(arg)
				if err != nil {
					results = append(results, reflect.Zero(outType))
					// 这个是 error
					results = append(results, reflect.ValueOf(err))
					return
				}
				// 你要在这里把调用信息拼凑起来
				// 服务名，方法名，参数值，参数类型不需要
				req := &Request{
					ServiceName: service.Name(),
					// 客户端和服务端可能叫不一样的名字
					// ServiceName: typ.PkgPath() + typ.Name(),
					// 服务名从哪里来？
					// 对应的是字段名
					MethodName: fieldType.Name,
					Data: bs,
				}
				// 要发送请求了啊
				// 要有一个接口

				// 不希望把 TCP 操作直接丢这里
				resp, err :=p.Invoke(ctx, req)

				if err != nil {
					results = append(results, reflect.Zero(outType))
					// 这个是 error
					results = append(results, reflect.ValueOf(err))
					return
				}
				// 第一个返回值，真的返回值，指向 GetIdResp
				first := reflect.New(outType).Interface()
				// 你要数据填进去 first
				// 这里就涉及到了序列化协议的问题
				// resp.Data => first 的转化

				// 我们现在先假定，这是用 JSON 来序列化的
				err = json.Unmarshal(resp.Data, first)

				results = append(results, reflect.ValueOf(first).Elem())

				// 这个是 error
				if err != nil {
					results = append(results, reflect.ValueOf(err))
				} else {
					results = append(results,  reflect.Zero(reflect.TypeOf(new(error)).Elem()))
				}



				// results = append(results,
				// 	reflect.New(fieldType.Type.Out(0)).Elem())
				// bs, err := json.Marshal(req)
				// if err != nil {
				// 	results = append(results, reflect.ValueOf(err))
				// 	return
				// }
				// // 你只是借助一下 error 来搞测试
				// results = append(results, reflect.ValueOf(errors.New(string(bs))))
				return
		})
		fieldValue.Set(fn)
	}
	return nil
}



type Service interface {
	Name() string
}