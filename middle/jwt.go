package middle

import (
	"github.com/gin-gonic/gin"
	"k8s/utils"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对登录接口放行
		if len(c.Request.URL.String()) >= 10 && c.Request.URL.String()[0:10] == "/api/login" || c.Request.URL.String()[0:9] == "/api/user" {
			c.Next()
		} else {
			// 获取Header中的Authorization
			token := c.Request.Header.Get("Authorization")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":  "请求未携带token，无权限访问",
					"data": nil,
				})
				c.Abort()
				return
			}
			// 解析token包含的信息
			claims, err := utils.ParseToken(token)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":  err.Error(),
					"data": nil,
				})
				c.Abort()
				return
			}
			c.Set("claims", claims)
			c.Next()
		}
	}
}
