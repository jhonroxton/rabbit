package routers

import (
	"github.com/gin-gonic/gin"
	"rabbit/api/v1"
)

func InitLeiSuDataRouterV1(g *gin.RouterGroup) {

	//鉴权接口
	group := g.Group("/api")
	{
		group.GET("/ping", v1.Ping)

	}

}
