package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type (
	// ModelInput represents a model input.
	ModelInput struct {
		Model         string              `json:"modelName"`     // Name of the model.
		InOrderFields []string            `json:"inOrderFields"` // Ordered fields.
		CSS           string              `json:"css"`           // Optional CSS, defaults to built in CSS.
		CardTemplates []CardTemplateInput `json:"cardTemplates"` // List of card templates.
	}

	// CardTemplateInput represents a card template input.
	CardTemplateInput struct {
		Name  string `json:"Name"`  // Card template name.
		Front string `json:"Front"` // Card front template.
		Back  string `json:"Back"`  // Card back template.
	}
)

// Models returns the list of model names for the current user.
func (c *Client) Models(ctx context.Context) ([]string, error) {
	request := ankiRequest{Action: "modelNames", Version: c.minVersion}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewReader(body))
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

// CreateModel creates a new model.
func (c *Client) CreateModel(ctx context.Context, model ModelInput) error {
	request := createModelRequest{
		ankiRequest: ankiRequest{Action: "createModel", Version: c.minVersion},
		Params:      model,
	}
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

type createModelRequest struct {
	ankiRequest
	Params ModelInput `json:"params"`
}
