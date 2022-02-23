package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			lg.Error("[Recovery from panic]",
				zap.Any("error", r),
				zap.String("\nrequest", c.Request.RequestURI+"\n"),
				zap.String("method", c.Request.Method+"\n"),
				zap.String("stack", string(debug.Stack())+"\n"),
			)
			c.JSON(http.StatusOK, gin.H{
				"code":   5000,
				"msg":    r,
				"result": nil,
			})
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	c.Next()
}

// recover错误，转string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
