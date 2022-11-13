package routers

import (
	"github.com/gin-gonic/gin"
	"rabbit/middleware"
)

func Routers() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	gin.SetMode(gin.DebugMode)

	//处理跨域
	r.Use(middleware.Cors)

	//统一添加路由组前缀 多服务器上线使用
	apiGroup := r.Group("")

	InitLeiSuDataRouterV1(apiGroup)

	return r

}
