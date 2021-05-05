package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type (
	StoreMediaInput struct {
		Filename string `json:"filename"`
		Data     string `json:"data"`
		Path     string `json:"path"`
		Url      string `json:"url"`
	}
)

func (c *Client) GetMedias(ctx context.Context, pattern string) ([]string, error) {
	request := getMediasRequest{
		ankiRequest: ankiRequest{Action: "getMediaFilesNames", Version: c.minVersion},
		Params:      getMediasParams{Pattern: pattern},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	var res []string
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) StoreMediaByData(ctx context.Context, media StoreMediaInput) (string, error) {
	request := storeMediaRequest{
		ankiRequest: ankiRequest{Action: "storeMediaFile", Version: c.minVersion},
		Params:      media,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req = req.WithContext(ctx)

	var res string
	if err := c.sendRequest(req, &res); err != nil {
		return "", err
	}

	return res, nil
}

type getMediasParams struct {
	Pattern string `json:"pattern"`
}

type getMediasRequest struct {
	ankiRequest
	Params getMediasParams `json:"params"`
}

type storeMediaRequest struct {
	ankiRequest
	Params StoreMediaInput `json:"params"`
}
