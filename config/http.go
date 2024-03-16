package config

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClientConfig struct {
	Client *http.Client
	Host string
}

func (httpClient HttpClientConfig) Post(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	url = fmt.Sprintf("%s%s", httpClient.Host, url)
	return http.NewRequestWithContext(ctx, "POST", url, body)
}

func (httpClient HttpClientConfig) Get(ctx context.Context, url string) (*http.Request, error) {
	url = fmt.Sprintf("%s%s", httpClient.Host, url)
	return http.NewRequestWithContext(ctx, "GET", url, nil)
}

func (httpClient HttpClientConfig) Put(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	url = fmt.Sprintf("%s%s", httpClient.Host, url)
	return http.NewRequestWithContext(ctx, "PUT", url, body)
}

func (httpClient HttpClientConfig) Delete(ctx context.Context, url string, body io.Reader) (*http.Request, error) {
	url = fmt.Sprintf("%s%s", httpClient.Host, url)
	return http.NewRequestWithContext(ctx, "DELETE", url, body)
}

func NewHttpContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10 * time.Second)
}