package errmgr

const (
	SUCCESS = 0

	Err_http_input_params_validate_error = 10001
	Err_http_input_params_empty          = 10002
	Err_http_imput_params_json_parse     = 10003

	Err_storage_root_not_found  = 20001
	Err_storage_path_invalid    = 20002
	Err_storage_list_failed     = 20003
	Err_storage_mkdir_failed    = 20004
	Err_storage_save_failed     = 20005
	Err_storage_read_failed     = 20006
	Err_storage_write_failed    = 20007
	Err_storage_lock_conflict   = 20008
	Err_storage_not_exists      = 20009
	Err_storage_ext_not_allowed = 20010
	Err_storage_file_too_large  = 20011

	Err_ai_forward_failed      = 30001
	Err_ai_forward_http_status = 30002
	Err_ai_forward_resp_parse  = 30003
)
