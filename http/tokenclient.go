package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/z26100/log-go"
	"io"
	"net/http"
)

type TokenClient struct {
	Url    string
	token  string
	client *http.Client
}

func NewTokenClient(url, token string, tr *http.Transport) *TokenClient {
	log.Infof("Set token client for url %s", url)
	return &TokenClient{
		Url:    url,
		token:  token,
		client: &http.Client{Transport: tr},
	}
}

func (c *TokenClient) SetToken(token string) {
	c.token = token
}

func (c TokenClient) Get(uri string, response interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.Url, uri), nil)
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, response)
}

func (c TokenClient) Post(uri string, requestBody interface{}, response interface{}) error {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.Url, uri), bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.client.Do(req)
	var jsonResponse []byte
	if resp != nil {
		jsonResponse, _ = io.ReadAll(resp.Body)
	}
	if err != nil {
		return errors.Join(err, errors.New(string(jsonResponse)))
	}
	return json.Unmarshal(jsonResponse, response)
}

func (c TokenClient) Delete(uri string, response interface{}) error {
	url := fmt.Sprintf("%s/%s", c.Url, uri)
	req, err := http.NewRequest("DELETE", url, nil)
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, response)
}
