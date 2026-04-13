module adm_bkd

go 1.24.0

toolchain go1.24.13

require (
	github.com/kataras/iris/v12 v12.2.11
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/spf13/viper v1.21.0
	go.uber.org/zap v1.27.1
	github.com/xuri/excelize/v2
)

require (
	github.com/go-sql-driver/mysql v1.9.3 // indirect: 保持与你参考项目一致（本项目不使用 DB）
	github.com/google/uuid v1.6.0 // indirect
)
