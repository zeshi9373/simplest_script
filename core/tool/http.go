package tool

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

var commonHttpClient *http.Client
var mx sync.Mutex

func NewHttp(url string, timeout time.Duration) *Http {
	mx.Lock()
	defer mx.Unlock()

	if commonHttpClient == nil {
		transport := &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 100,
		}
		commonHttpClient = &http.Client{Transport: transport}
	}

	return &Http{
		Url:        url,
		Timeout:    timeout,
		httpClient: commonHttpClient,
	}
}

func (h *Http) PostByForm(header map[string]string, data map[string]string) ([]byte, error) {
	urlV := url.Values{}
	for k, v := range data {
		urlV.Add(k, v)
	}

	ctx := context.Background()
	if h.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.Timeout)
		defer cancel()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, h.Url, strings.NewReader(urlV.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, value := range header {
		request.Header.Set(key, value)
	}

	return h.httpDo(request)
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

	request.Header.Set("Content-Type", "application/json")

	return h.httpDo(request)
}

// 发送HTTP GET请求
func (h *Http) Get(header map[string]string, data map[string]string) ([]byte, error) {
	ctx := context.Background()

	if h.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, h.Timeout)
		defer cancel()
	}

	fullUrl := h.Url

	if len(data) > 0 {
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
					if strings.Contains(kv[1], "%") {
						unescape, err := url.QueryUnescape(kv[1])

						if err != nil || len(unescape) == 0 {
							params.Add(kv[0], kv[1])
						} else {
							params.Add(kv[0], unescape)
						}
					} else {
						params.Add(kv[0], kv[1])
					}
				}
			}
		}

		for key, value := range data {
			params.Set(key, value)
		}

		rawUrl.RawQuery = params.Encode()
		fullUrl = rawUrl.String()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	return h.httpDo(request)
}

// 发送HTTP请求
func (h *Http) httpDo(request *http.Request) ([]byte, error) {
	response, err := h.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}
