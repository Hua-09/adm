package main

import (
	"adm_bkd/api"
	"adm_bkd/config"
	"log"

	"github.com/kataras/iris/v12"
	requestLogger "github.com/kataras/iris/v12/middleware/logger"
)

func main() {
	if err := config.LoadConfig("./config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	app := iris.New()

	// 简化：先使用 Iris 自带 logger 中间件；如果你后续把参考项目 logs 包拷过来，可在这里替换
	app.Use(requestLogger.New())
	app.Logger().SetLevel(config.GlobalConfig.GetLogLevel())

	apiApp := app.Party("/api")

	// health
	apiApp.Get("/ping", api.HealthCheck)

	// teach repo APIs
	apiApp.Post("/teach/dept/list", api.TeachDeptList)
	apiApp.Post("/teach/date/list", api.TeachDateList)
	apiApp.Post("/teach/teacher/list", api.TeachTeacherList)
	apiApp.Post("/teach/file/list", api.TeachFileList)

	apiApp.Post("/teach/file/upload", api.TeachFileUpload)

	apiApp.Post("/teach/analyze/create", api.TeachAnalyzeCreate)
	apiApp.Post("/teach/analyze/status", api.TeachAnalyzeStatus)

	apiApp.Post("/teach/result/md/get", api.TeachResultMdGet)
	apiApp.Post("/teach/result/json/get", api.TeachResultJsonGet)
	apiApp.Post("/teach/result/viz/get", api.TeachResultVizGet)

	addr := config.GlobalConfig.GetServerAddr()
	app.Run(iris.Addr(addr))
}
