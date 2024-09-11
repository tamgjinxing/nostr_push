package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPClientConfig 是 HTTP 客户端配置
type HTTPClientConfig struct {
	Timeout         time.Duration // 请求超时时间
	MaxIdleConns    int           // 最大空闲连接数
	IdleConnTimeout time.Duration // 空闲连接的超时时间
}

// HTTPClient 工具类结构体，包含 http.Client 和 配置信息
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建并返回一个带连接池的 HTTPClient 工具类
func NewHTTPClient(config HTTPClientConfig) *HTTPClient {
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,    // 最大空闲连接数
		IdleConnTimeout:     config.IdleConnTimeout, // 空闲连接超时时间
		MaxIdleConnsPerHost: config.MaxIdleConns,    // 每个主机的最大空闲连接数
	}

	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}

	return &HTTPClient{client: httpClient}
}

// Get 发送 GET 请求
func (h *HTTPClient) Get(url string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

// Post 发送 POST 请求（JSON 格式）
func (h *HTTPClient) Post(url string, headers map[string]string, data interface{}) ([]byte, int, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

// PostForm 发送 POST 请求（表单形式）
func (h *HTTPClient) PostForm(url string, headers map[string]string, formData map[string]string) ([]byte, int, error) {
	form := ""
	for key, value := range formData {
		form += fmt.Sprintf("%s=%s&", key, value)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(form))
	if err != nil {
		return nil, 0, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

// Put 发送 PUT 请求
func (h *HTTPClient) Put(url string, headers map[string]string, data interface{}) ([]byte, int, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

// Delete 发送 DELETE 请求
func (h *HTTPClient) Delete(url string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, 0, err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

// CheckStatus 检查 HTTP 响应状态码
func CheckStatus(statusCode int) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}
	return errors.New(fmt.Sprintf("unexpected status code: %d", statusCode))
}
