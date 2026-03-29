package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (c *Client) httpGetOnce(path string, out any) error {
	time.Sleep(650 * time.Millisecond) // limit = 100 request/minute → 0.6 sec/request
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}
	return json.Unmarshal(body, out)
}

func (c *Client) HttpGet(path string, out any) error {
	for attempt := 1; attempt <= httpGetRetries; attempt++ {
		err := c.httpGetOnce(path, out)
		if err == nil {
			break
		}
		if attempt == httpGetRetries {
			return err
		}
		log.Warn("httpGet failed, retrying", "attempt", attempt, "err", err)
		time.Sleep(httpGetRetryDelay)
	}
	return nil
}
