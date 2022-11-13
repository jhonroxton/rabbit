package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//配置中间件 跨域

func Cors(c *gin.Context) {

	// gin设置响应头，设置跨域
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin") //请求头部
	if origin != "" {
		//接收客户端发送的origin （重要！）
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		//设置所有的请求都允许
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		//服务器支持的所有跨域请求的方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		//允许跨域设置可以返回其他子段，可以自定义字段
		c.Header("Access-Control-Allow-Headers", "Content-Length, X-CSRF-Token, Token,session")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,Authorization, Token,X-Token,X-User-Id,x-token")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,charset=utf-8")
		//设置缓存时间
		c.Header("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		c.Header("Access-Control-Allow-Credentials", "true")

		//这句我也不知道是干嘛的`反正vue前端需要跨域就的设置这句
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With , yourHeaderFeild")
	}

	//允许类型校验
	if method == "OPTIONS" {
		c.JSON(http.StatusOK, "ok!")
	}
	c.Next()
}
