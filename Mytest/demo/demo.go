package demo

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Start() {
	s := gin.Default()
	fmt.Println("start")
	s.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	s.Run(":8080")
}
