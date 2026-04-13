package api_result

// Result is the standard JSON response envelope.
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK returns a successful result.
func OK(data interface{}) Result {
	return Result{Code: 0, Message: "ok", Data: data}
}

// Fail returns an error result.
func Fail(code int, message string) Result {
	return Result{Code: code, Message: message}
}
