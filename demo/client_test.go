package demo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitClientProxy(t *testing.T) {
	testCases := []struct{
		name string

		service *UserService

		wantErr error
	}{
		{
			name: "user service",
			service: &UserService{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := InitClientProxy(tc.service)
			assert.Equal(t, tc.wantErr, err)
			// 缺乏校验手段
			_, _ = tc.service.GetById(context.Background(), &GetByIdReq{})
		})
	}
}

type UserService struct {
	GetById func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
}

func (u *UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {

}

type GetByIdResp struct {

}

// type UserServer struct {
// 	GetByIdA func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
// }
//
//
// type UserClient struct {
// 	GetByIdB func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
// }