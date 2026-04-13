package err_mgr

var errMessages = map[int]string{
	ErrOK:            "ok",
	ErrInvalidParam:  "invalid parameter",
	ErrStorageRead:   "storage read error",
	ErrStorageWrite:  "storage write error",
	ErrStorageDelete: "storage delete error",
	ErrAIForward:     "AI forwarding error",
	ErrInternal:      "internal server error",
}

// ErrStr returns the human-readable message for a given error code.
func ErrStr(code int) string {
	if msg, ok := errMessages[code]; ok {
		return msg
	}
	return "unknown error"
}
