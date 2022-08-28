package gin

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestUserController_GetUser(t *testing.T) {
	g := gin.Default()
	ctrl := &UserController{}
	g.GET("/user/*", ctrl.GetUser)
	_ = g.Run(":8082")
}
