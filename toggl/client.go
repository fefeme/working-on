package toggl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BaseURL       = "https://api.track.toggl.com/api/v8"
	shortDuration = 1000 * time.Millisecond
)

type Toggl struct {
	TimeEntries     *TimeEntries
	TaskClient      *TaskClient
	WorkspaceClient *WorkspaceClient
}

type Client struct {
	BaseURL    string
	apiToken   string
	HttpClient *http.Client
}

type Message struct {
	endpoint string
	method   string
	payload  *bytes.Buffer
}

func NewToggl(apiToken string) *Toggl {
	client := Client{
		BaseURL:  BaseURL,
		apiToken: apiToken,
		HttpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
	return &Toggl{
		TimeEntries: &TimeEntries{
			client: &client,
		},
		TaskClient: &TaskClient{
			client: &client,
		},
		WorkspaceClient: &WorkspaceClient{
			client: &client,
		},
	}
}

func (c *Client) NewMessage(method string, endpoint string, data interface{}) (*Message, error) {
	payload := new(bytes.Buffer)
	if data != nil {
		err := json.NewEncoder(payload).Encode(data)
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		endpoint: fmt.Sprintf("%s/%s", c.BaseURL, endpoint),
		method:   method,
		payload:  payload,
	}, nil

}

func (c *Client) SendRequest(message *Message) (*json.RawMessage, error) {
	var (
		req *http.Request
		err error
	)

	if message.payload != nil {
		req, err = http.NewRequest(message.method, message.endpoint, message.payload)
	} else {
		req, err = http.NewRequest(message.method, message.endpoint, nil)
	}

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
	defer cancel()

	req = req.WithContext(ctx)
	req.SetBasicAuth(c.apiToken, "api_token")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		return nil, errors.New("Rate limit hit.")
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println("Return Value")
		fmt.Println(res.Status)
		fmt.Println(b)
		return nil, fmt.Errorf("%s, status code: %d", b, res.StatusCode)
	}

	js, err := ioutil.ReadAll(res.Body)

	var raw json.RawMessage
	if json.Unmarshal(js, &raw) != nil {
		return nil, err
	}

	//fmt.Printf("Raw Json %s\n", js)

	return &raw, err

}
