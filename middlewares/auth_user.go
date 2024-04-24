package middlewares

import (
	"OJ/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 验证用户
func AuthUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO:Check if user is admin
		auth := c.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Unauthorized Authorization",
			})
			return
		}
		if userClaim == nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Unauthorized Admin",
			})
			return
		}
		// 往user里设置userClaim
		c.Set("user", userClaim)
		// 去执行下一个中间件
		c.Next()
	}
}
