package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// Exit schedules a request to gracefully close Anki.
// The operation is asynchronous, so it will return immediately
// and won't wait until the Anki process actually terminates.
func (c *Client) Exit(ctx context.Context) error {
	request := ankiRequest{Action: "guiExitAnki", Version: c.minVersion}
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewReader(body))
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
