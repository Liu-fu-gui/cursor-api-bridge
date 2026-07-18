package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ModelListRequest struct {
	BaseURL string `json:"baseURL"`
	APIKey  string `json:"apiKey"`
}

type modelListResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

// ListOpenAIModels 从 OpenAI-compatible /v1/models 接口读取模型列表。
func (s *ProxyService) ListOpenAIModels(input ModelListRequest) ([]string, error) {
	if s == nil || s.publicClient == nil {
		return nil, errors.New("网络服务尚未初始化")
	}
	baseText := strings.TrimSpace(input.BaseURL)
	apiKey := strings.TrimSpace(input.APIKey)
	if baseText == "" || apiKey == "" {
		return nil, errors.New("URL 和 API Key 不能为空")
	}
	base, err := url.Parse(baseText)
	if err != nil || base.Scheme == "" || base.Host == "" {
		return nil, errors.New("API URL 格式不正确")
	}
	path := strings.TrimRight(base.Path, "/")
	if strings.HasSuffix(path, "/v1") {
		path += "/models"
	} else {
		path += "/v1/models"
	}
	base.Path = path
	base.RawQuery = ""
	base.Fragment = ""

	req, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")
	resp, err := s.publicClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("连接模型服务失败: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		detail := strings.TrimSpace(string(body))
		if len(detail) > 500 {
			detail = detail[:500]
		}
		return nil, fmt.Errorf("模型服务返回 %d: %s", resp.StatusCode, detail)
	}
	var payload modelListResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, errors.New("模型列表不是有效的 OpenAI 格式")
	}
	seen := make(map[string]struct{}, len(payload.Data))
	models := make([]string, 0, len(payload.Data))
	for _, item := range payload.Data {
		id := strings.TrimSpace(item.ID)
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		models = append(models, id)
	}
	if len(models) == 0 {
		return nil, errors.New("接口连接成功，但没有返回模型")
	}
	return models, nil
}
