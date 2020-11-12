package number

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{
		url: url,
	}
}

func (c *Client) Info(num int) (string, error) {
	numStr := fmt.Sprintf("/%d", num)
	resp, err := http.Get(c.url + numStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
