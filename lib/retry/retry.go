package retry

import (
	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"io"
	"time"
)

// 客户端 重试
func DoRetry(request func() error, times int, delay time.Duration) error {
	return retry.Do(request, retry.Attempts(uint(times)), retry.Delay(delay), retry.DelayType(retry.BackOffDelay))
}

// 客户端 http
type Client struct {
	path   string
	Client *resty.Request
}

func ClientBuild() *Client {
	client := resty.New()
	return &Client{Client: client.R()}
}
func (this *Client) SetQueryParams(params map[string]string) *Client {
	this.Client.SetQueryParams(params)
	return this
}
func (this *Client) SetFormData(postData map[string]string) *Client {
	this.Client.SetFormData(postData)
	return this
}
func (this *Client) SetHeader(key string, value string) *Client {
	this.Client.SetHeader(key, value)
	return this
}

func (this *Client) Url(key string) *Client {
	this.path = key
	return this
}
func (this *Client) Json(key string) *Client {
	this.SetHeader("Content-Type", "application/json")
	this.SetHeader("Accept", "application/json")
	return this
}
func (this *Client) SetBody(data interface{}) *Client {
	this.Client.SetBody(data)
	return this
}
func (this *Client) SetFileReader(param, fileName string, reader io.Reader) *Client {
	this.Client.SetFileReader(param, fileName, reader)
	return this
}
func (this *Client) SetFile(param, filePath string) *Client {
	this.Client.SetFile(param, filePath)
	return this
}
func (this *Client) Get() (*resty.Response, error) {
	return this.Client.Get(this.path)
}
func (this *Client) Post() (*resty.Response, error) {
	return this.Client.Post(this.path)
}
func (this *Client) Put() (*resty.Response, error) {
	return this.Client.Put(this.path)
}
func (this *Client) Delete() (*resty.Response, error) {
	return this.Client.Delete(this.path)
}
func (this *Client) Head() (*resty.Response, error) {
	return this.Client.Head(this.path)
}
