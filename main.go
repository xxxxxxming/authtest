package main

import (
	"main/middlewares"
	"main/routers"
	"main/utils"
	"main/utils/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 读取配置文件
	cfg, err := config.ParseConfig("./configs/app.json")
	if err != nil {
		panic(utils.Errors(err, "read config file error"))
	}
	// 初始化日志记录器
	if err := middlewares.InitLogger(&cfg.LogConfig); err != nil {
		panic(utils.Errors(err, "init logger failed, err"))
	}
	app := gin.Default()
	// 设置跨域
	app.Use(middlewares.Cors())
	// 设置jwt认证和权限校验
	app.Use(middlewares.Authorization())
	// 设置异常捕获中间件
	app.Use(middlewares.Recover)
	// 设置日志中间件
	app.Use(middlewares.GinLogger())

	utils.InitTree()
	utils.InitAuth()
	// 注册路由
	registerRoutes(app)

	zap.L().Info("系统启动成功")

	app.Run(cfg.App.Apphost + ":" + cfg.App.Appport)

}

func registerRoutes(c *gin.Engine) {
	// routes.TestRouter(c)
	routers.TestRouter(c)
}
