package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BasicAuthClient struct {
	Url      string
	Username string
	Password string
	client   *http.Client
}

func (c BasicAuthClient) Post(uri string, requestBody interface{}, response interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c BasicAuthClient) Delete(uri string, response interface{}) error {
	//TODO implement me
	panic("implement me")
}

func NewBasicAuthClient(url, username, password string, tr *http.Transport) *BasicAuthClient {
	return &BasicAuthClient{
		Url:      url,
		Username: username,
		Password: password,
		client:   &http.Client{Transport: tr},
	}
}

func (c BasicAuthClient) Get(uri string, response interface{}) error {
	url := fmt.Sprintf("%s/%s", c.Url, uri)
	req, err := http.NewRequest("GET", url, nil)
	//req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)
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
