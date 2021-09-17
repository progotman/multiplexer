package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Requester struct {
	Client http.Client
}

func (d Requester) GetContent(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := d.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("cannot read body: %w", err)
		}

		return "", Error{Code: resp.StatusCode, Message: string(bodyBytes)}
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body: %w", err)
	}

	return string(bodyBytes), nil
}

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}
