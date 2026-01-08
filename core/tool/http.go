package tool

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Http struct {
	Url             string
	Header          map[string]string
	Body            []byte
	Timeout         time.Duration
	httpClient      *http.Client
	ResponseHeaders map[string]string
}

func NewHttp(url string, timeout time.Duration) *Http {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConnsPerHost: 100,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	return &Http{
		Url:        url,
		Timeout:    timeout,
		httpClient: httpClient,
	}
}

func (h *Http) PostByForm(header map[string]string, data map[string]string) ([]byte, error) {
	urlV := url.Values{}

	for k, v := range data {
		urlV.Add(k, v)
	}

	// 创建一个 POST 请求
	resp, err := http.PostForm(h.Url, urlV)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	return body, nil
}

// 发送HTTP POST请求
func (h *Http) Post(header map[string]string, data []byte) ([]byte, error) {
	ctx := context.Background()

	if h.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.Timeout)
		defer cancel()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, h.Url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	//设置请求头Content-Type
	request.Header.Set("Content-Type", "application/json")

	result, err := h.httpDo(request)

	if err != nil {
		return nil, err
	}

	return result, err
}

// 发送HTTP GET请求
func (h *Http) Get(header map[string]string, data map[string]string) ([]byte, error) {
	ctx := context.Background()

	if h.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.Timeout)
		defer cancel()
	}

	params := url.Values{}
	rawUrl, err := url.Parse(h.Url)

	if err != nil {
		return nil, err
	}

	if len(rawUrl.RawQuery) > 0 {
		args := strings.Split(rawUrl.RawQuery, "&")

		for _, arg := range args {
			kv := strings.Split(arg, "=")
			if len(kv) == 2 {
				params.Add(kv[0], kv[1])
			}
		}
	}

	for key, value := range data {
		params.Set(key, value)
	}

	rawUrl.RawQuery = params.Encode()
	url := rawUrl.String()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	result, err := h.httpDo(request)

	if err != nil {
		return nil, err
	}

	return result, err
}

// 发送HTTP请求
func (h *Http) httpDo(request *http.Request) ([]byte, error) {
	h.httpClient.Timeout = time.Duration(h.Timeout)
	response, err := h.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
