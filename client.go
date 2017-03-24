package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	apiVersion      = "v3"
	acceptHeader    = "application/vnd.github." + apiVersion + "+json"
	userAgentHeader = "boss/" + "1.0"
	apiBase         = "https://api.github.com/"
)

type Client struct {
	Client    *http.Client
	APIBase   *url.URL
	Token     string
	UserAgent string
}

func NewClient(token string, c *http.Client) (*Client, error) {

	if c == nil {
		c = http.DefaultClient
	}

	APIBase, err := url.Parse(apiBase)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:    c,
		APIBase:   APIBase,
		Token:     token,
		UserAgent: userAgentHeader,
	}, nil
}

func (self *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := self.APIBase.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", acceptHeader)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", self.Token))
	req.Header.Add("User-Agent", self.UserAgent)

	return req, nil
}

func (self *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := self.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil
			}
		}
	}
	return resp, nil
}
