package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// Version gets the exposed version of AnkiConnect's API.
func (c *Client) Version(ctx context.Context) (int, error) {
	request := ankiRequest{Action: "version", Version: c.minVersion}
	body, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)

	var res int
	if err := c.sendRequest(req, &res); err != nil {
		return 0, err
	}

	return res, nil
}

// Sync synchronizes the local Anki collections with AnkiWeb.
func (c *Client) Sync(ctx context.Context) error {
	request := ankiRequest{Action: "sync", Version: c.minVersion}
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)

	var res interface{}
	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}
