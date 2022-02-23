package middlewares

import (
	"main/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := strings.ToLower(c.Request.Method)
		// path := c.Request.URL.Path
		parmas := c.Request.Header
		var code uint
		var msg string
		var errflag = false
		// 将 fullpath,method和query传入解析树函数中
		str := utils.Root.ParseUrlTree(c.FullPath(), method, c.Request.URL.RawQuery)
		if str != "" {
			// 获取需要验证的函数对应的字符串
			tokenHandle, b := utils.AuthMap[str]
			if b {
				token := parmas["Authorization"]
				if token == nil {
					code = 4001
					msg = "缺失Authorization请求头"
					errflag = true
				} else {
					claims, err := tokenHandle.TokenAuth(token[0])
					if err != nil {
						code = 4003
						msg = err.Error()
						errflag = true
					} else if claims == nil {
						code = 4003
						msg = "token校验失败"
						errflag = true
					}
				}

			}
		}
		if errflag == true {
			c.JSON(200, gin.H{
				"code":   code,
				"msg":    msg,
				"result": nil,
			})
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
		c.Next()
	}
}
