package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sfmunoz/clickup-archive/internal/api"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

const (
	baseURL           = "https://api.clickup.com/api/v2"
	httpGetRetries    = 5
	httpGetRetryDelay = time.Second
)

type Client struct {
	token string
}

func NewClient() (*Client, error) {
	token := os.Getenv("CLICKUP_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("CLICKUP_TOKEN env var is required")
	}
	return &Client{token: token}, nil
}

func (c *Client) httpGetOnce(url string) ([]byte, error) {
	time.Sleep(650 * time.Millisecond)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}
	return body, nil
}

func (c *Client) HttpGet(url string) ([]byte, error) {
	var data []byte
	for attempt := 1; attempt <= httpGetRetries; attempt++ {
		var err error
		data, err = c.httpGetOnce(url)
		if err == nil {
			break
		}
		if attempt == httpGetRetries {
			return nil, err
		}
		log.Warn("httpGetBytes failed, retrying", "attempt", attempt, "err", err)
		time.Sleep(httpGetRetryDelay)
	}
	return data, nil
}

func (c *Client) HttpGetJson(path string, out any) error {
	body, err := c.HttpGet(baseURL + path)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func (c *Client) GetTask(taskID string) (*api.Task, error) {
	var task api.Task
	if err := c.HttpGetJson("/task/"+taskID, &task); err != nil {
		return nil, err
	}
	return &task, nil
}
