package apiresult

import (
	errmgr "adm_bkd/utils/err_mgr"
)

type APIResult struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Msg     string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewAPIResult(code int, data interface{}) APIResult {
	msg := errmgr.ErrStr(code)
	ret := APIResult{false, code, msg, data}
	if code == 0 {
		ret.Success = true
	}
	return ret
}

func NewAPIResultWithStr(code int, msg string, data interface{}) APIResult {
	ret := APIResult{false, code, msg, data}
	if code == 0 {
		ret.Success = true
	}
	return ret
}
