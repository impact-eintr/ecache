package client

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type httpClient struct {
	*http.Client
	server string
}

func newHttpClient(server string) *httpClient {
	client := &http.Client{Transport: &http.Transport{
		MaxIdleConnsPerHost: 1,
	}}
	return &httpClient{client, "http://" + server + "/cache/"}
}

func (c *httpClient) get(key string) ([]byte, error) {
	resp, err := c.Get(c.server + key)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("cache miss!")
	}

	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	return nil, errors.New(resp.Status)

}

func (c *httpClient) set(key string, value []byte) error {
	req, err := http.NewRequest(
		http.MethodPut, c.server+key, bytes.NewReader(value))

	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil

}

func (c *httpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		cmd.Value, cmd.Error = c.get(cmd.Key)
	}
}

func (c *httpClient) PipelinedRun([]*Cmd) {}
