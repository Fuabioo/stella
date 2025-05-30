package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"resty.dev/v3"
)

type ProgressReport struct {
	Total   uint   `json:"total"`
	Current uint   `json:"current"`
	Text    string `json:"text"`
	Loading bool   `json:"loading"`
}

type Client struct {
	*resty.Client
	step             time.Duration
	backoff          time.Duration
	failureThreshold uint
	running          bool
}

func NewClient(uri string,
	step,
	backoff time.Duration,
	failureThreshold uint,
	superDebug bool,
) *Client {
	return &Client{
		Client: resty.New().
			SetLogger(log.Default()).
			SetDebug(superDebug).
			SetBaseURL(uri),
		step:             step,
		backoff:          backoff,
		failureThreshold: failureThreshold,
	}
}

func (client *Client) GetProgress() (*ProgressReport, error) {
	if client == nil || client.Client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	var result ProgressReport

	resp, err := client.R().
		SetResult(&result).
		Get("")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch progress: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to fetch progress: statusCode=%s", resp.Status())
	}

	return &result, nil
}

func (client *Client) Observe(callback func(*ProgressReport, error)) {
	if client == nil || client.running {
		return
	}
	go func() {
		for {
			time.Sleep(client.step)
			result, err := client.GetProgress()
			callback(result, err)
			if err != nil {
				time.Sleep(client.backoff)
			}
		}
	}()

}
