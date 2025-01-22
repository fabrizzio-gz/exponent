package exponent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Client struct {
	cnf *Config
}

func NewClient(opts ...Option) *Client {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	withDefaults(c)
	return &Client{c}
}

// Publish sends a single push notification
// @param push_message: A PushMessage object
// @return an array of PushResponse objects which contains the results.
// @return error if any requests failed
func (c *Client) PublishSingle(ctx context.Context, msg *Message) (*MessageResponse, error) {
	responses, err := c.publish(ctx, []*Message{msg})
	if err != nil {
		return nil, err
	}
	return responses[0], nil
}

// PublishMultiple sends multiple push notifications at once
// @param push_messages: An array of PushMessage objects.
// @return an array of PushResponse objects which contains the results.
// @return error if the request failed
func (c *Client) Publish(ctx context.Context, msgs []*Message) ([]*MessageResponse, error) {
	return c.publish(ctx, msgs)
}

func (c *Client) publish(ctx context.Context, msgs []*Message) ([]*MessageResponse, error) {
	// Validate the messages
	for _, message := range msgs {
		if len(message.To) == 0 {
			return nil, errors.New("no recipients")
		}
		for _, recipient := range message.To {
			if recipient == nil || *recipient == "" {
				return nil, errors.New("invalid push token")
			}
		}
	}
	url := fmt.Sprintf("%s%s/push/send", c.cnf.Host, c.cnf.ApiURL)
	jsonBytes, err := json.Marshal(msgs)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if c.cnf.AcessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.cnf.AcessToken)
	}
	resp, err := c.cnf.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err = checkStatus(resp); err != nil {
		return nil, err
	}
	var r *Response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	if r.Errors != nil {
		return nil, errors.New("invalid request")
	}
	if r.Data == nil {
		return nil, NewServerError("invalid server response", resp, r, nil)
	}
	// Its closed because MessageItem has own Status, Message, etc.
	// So, we don't need to check the length of r.Data
	// if len(msgs) != len(r.Data) {
	// 	errMsg := fmt.Sprintf("mismatched response length. Expected %d receipts but only received %d", len(msgs), len(r.Data))
	// 	return nil, NewServerError(errMsg, resp, r, nil)
	// }
	for i := range r.Data {
		r.Data[i].MessageItem = msgs[i]
	}
	return r.Data, nil
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= 299 {
		return nil
	}
	return fmt.Errorf("invalid response (%d %s)", resp.StatusCode, resp.Status)
}
