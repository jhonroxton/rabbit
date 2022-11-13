package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	fmt.Println("请求成功")
}
