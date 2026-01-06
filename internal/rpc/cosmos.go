package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type StakingValidatorResponse struct {
	Validator struct {
		OperatorAddress string `json:"operator_address"`
		Jailed          bool   `json:"jailed"`
		Status          string `json:"status"`
		Tokens          string `json:"tokens"`
	} `json:"validator"`
}

type CosmosClient struct {
	baseURL string
	client  *http.Client
}

func NewCosmos(baseURL string) *CosmosClient {
	return &CosmosClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func NewCosmosWithTimeout(baseURL string, timeout time.Duration) *CosmosClient {
	return &CosmosClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *CosmosClient) Validator(valoper string) (*StakingValidatorResponse, error) {
	url := fmt.Sprintf(
		"%s/cosmos/staking/v1beta1/validators/%s",
		c.baseURL,
		valoper,
	)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var out StakingValidatorResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}
