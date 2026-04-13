package storage

import (
	apiresult "adm_bkd/utils/api_result"
	errmgr "adm_bkd/utils/err_mgr"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// ForwardToAiServer 将 payload 转发到 Python AI 服务，返回响应 body
// 说明：你的 Python 服务返回格式可能不同；这里返回原始 bytes，由上层做兼容解析
func ForwardToAiServer(apiUrl string, payload map[string]interface{}) ([]byte, int) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, errmgr.Err_ai_forward_failed
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errmgr.Err_ai_forward_failed
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errmgr.Err_ai_forward_failed
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return b, errmgr.Err_ai_forward_http_status
	}

	// 兼容：如果 Python 返回了统一包裹（success/code/message/data），也让上层自己解析
	_ = apiresult.APIResult{} // 引入以保持与你项目风格一致（非必须）

	return b, errmgr.SUCCESS
}
