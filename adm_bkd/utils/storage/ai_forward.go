package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AIRequest is the payload sent to the AI service.
type AIRequest struct {
	DocID   string `json:"doc_id"`
	RawText string `json:"raw_text"`
}

// AIResponse is the payload received from the AI service.
type AIResponse struct {
	DocID  string `json:"doc_id"`
	Result string `json:"result"`
}

// ForwardToAI sends parsed document content to the configured AI endpoint and
// returns the AI service response.
func ForwardToAI(endpoint string, timeoutSeconds int, docID, rawText string) (AIResponse, error) {
	payload := AIRequest{DocID: docID, RawText: rawText}
	body, err := json.Marshal(payload)
	if err != nil {
		return AIResponse{}, fmt.Errorf("marshal AI request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return AIResponse{}, fmt.Errorf("build AI request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return AIResponse{}, fmt.Errorf("call AI endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AIResponse{}, fmt.Errorf("AI endpoint returned %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return AIResponse{}, fmt.Errorf("read AI response: %w", err)
	}

	var aiResp AIResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return AIResponse{}, fmt.Errorf("unmarshal AI response: %w", err)
	}
	return aiResp, nil
}
