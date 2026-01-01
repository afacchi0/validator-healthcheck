package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type StatusResponse struct {
	Result struct {
		NodeInfo struct {
			ID      string `json:"id"`
			Network string `json:"network"`
		} `json:"node_info"`
		SyncInfo struct {
			CatchingUp bool `json:"catching_up"`
		} `json:"sync_info"`
	} `json:"result"`
}

type ValidatorsResponse struct {
	Result struct {
		Validators []struct {
			Address          string `json:"address"`
			VotingPower      string `json:"voting_power"`
			ProposerPriority string `json:"proposer_priority"`
		} `json:"validators"`
	} `json:"result"`
}

type Client struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Status() (*StatusResponse, error) {
	url := fmt.Sprintf("%s/status", c.baseURL)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var out StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) Validators() (*ValidatorsResponse, error) {
	url := fmt.Sprintf("%s/validators", c.baseURL)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var out ValidatorsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
