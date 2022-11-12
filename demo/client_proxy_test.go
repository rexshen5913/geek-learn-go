package demo

import "testing"

func TestClient_Invoke(t *testing.T) {
	cli := &Client{}
	InitClientProxy(&UserServiceClient{}, cli)
}
