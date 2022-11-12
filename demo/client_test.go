package demo

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitClientProxy(t *testing.T) {
	testCases := []struct{
		name string

		service *UserServiceClient
		p *mockProxy

		wantReq *Request
		wantInitErr error
		wantErr error
		wantResp *GetByIdResp
	}{
		{
			name: "user service",
			p: &mockProxy{
				result: []byte(`{"name": "Tom"}`),
			},
			service: &UserServiceClient{},
			wantReq:  &Request{
				ServiceName: "user-service",
				MethodName: "GetById",
				Data: []byte(`{"id": 13}`),
			},
			wantResp: &GetByIdResp{
				Name: "Tom",
			},
		},

		{
			name: "proxy return error",
			p: &mockProxy{
				err: errors.New("mock error"),
			},
			service: &UserServiceClient{},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitClientProxy(tc.service, tc.p)
			assert.Equal(t, tc.wantInitErr, err)
			if err != nil {
				return
			}
			// 缺乏校验手段
			resp, err := tc.service.GetById(context.Background(), &GetByIdReq{Id: 13})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			// 断言 p 的数据
			assert.Equal(t, tc.wantReq, tc.p.req)
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

// 可以考虑使用 gomock
type mockProxy struct {
	req *Request
	err error
	result []byte
}

func (m *mockProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
	m.req = req
	return &Response{
		Data: m.result,
	}, m.err
}

type UserServiceClient struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserServiceClient) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int64
}

type GetByIdResp struct {
	Name string `json:"name"`
}

// type UserServer struct {
// 	GetByIdA func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
// }
//
//
// type UserClient struct {
// 	GetByIdB func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
// }