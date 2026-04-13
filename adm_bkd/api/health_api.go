package api

import (
	apiresult "adm_bkd/utils/api_result"
	errmgr "adm_bkd/utils/err_mgr"

	"github.com/kataras/iris/v12"
)

func HealthCheck(ctx iris.Context) {
	ctx.JSON(apiresult.NewAPIResult(errmgr.SUCCESS, "pong"))
}
