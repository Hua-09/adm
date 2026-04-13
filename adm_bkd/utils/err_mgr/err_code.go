package errmgr

// ErrStr 返回错误码对应的文本说明（用于 APIResult.message）
func ErrStr(code int) string {
	switch code {
	case SUCCESS:
		return "success"

	case Err_http_input_params_validate_error:
		return "input params validate error"
	case Err_http_input_params_empty:
		return "input params empty"
	case Err_http_imput_params_json_parse:
		return "input params json parse failed"

	case Err_storage_root_not_found:
		return "storage root dir not found"
	case Err_storage_path_invalid:
		return "storage path invalid"
	case Err_storage_list_failed:
		return "storage list failed"
	case Err_storage_mkdir_failed:
		return "storage mkdir failed"
	case Err_storage_save_failed:
		return "storage save file failed"
	case Err_storage_read_failed:
		return "storage read failed"
	case Err_storage_write_failed:
		return "storage write failed"
	case Err_storage_lock_conflict:
		return "analyze task already running"
	case Err_storage_not_exists:
		return "target not exists"
	case Err_storage_ext_not_allowed:
		return "file ext not allowed"
	case Err_storage_file_too_large:
		return "file too large"

	case Err_ai_forward_failed:
		return "ai forward failed"
	case Err_ai_forward_http_status:
		return "ai service http status not ok"
	case Err_ai_forward_resp_parse:
		return "ai service response parse failed"

	default:
		return "unknown error"
	}
}
