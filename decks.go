package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// Decks gets the complete list of deck names for the current user.
func (c *Client) Decks(ctx context.Context) ([]string, error) {
	request := ankiRequest{Action: "deckNames", Version: c.minVersion}
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

// CreateDeck creates a new empty deck.
// Will not overwrite a deck that exists with the same name.
func (c *Client) CreateDeck(ctx context.Context, name string) (int, error) {
	request := createDeckRequest{
		ankiRequest: ankiRequest{Action: "createDeck", Version: c.minVersion},
		Params:      createDeckParams{Deck: name},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)

	var res int
	if err := c.sendRequest(req, &res); err != nil {
		return 0, nil
	}

	return res, nil
}

// DeleteDecks deletes decks with the given names.
func (c *Client) DeleteDecks(ctx context.Context, names []string) error {
	// TODO: cards too option?
	request := deleteDecksRequest{
		ankiRequest: ankiRequest{Action: "deleteDecks", Version: c.minVersion},
		Params:      deleteDecksParams{Decks: names},
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

type createDeckParams struct {
	Deck string `json:"deck"`
}

type createDeckRequest struct {
	ankiRequest
	Params createDeckParams `json:"params"`
}

type deleteDecksParams struct {
	Decks    []string `json:"decks"`
	CardsToo bool     `json:"cardsToo"`
}

type deleteDecksRequest struct {
	ankiRequest
	Params deleteDecksParams `json:"params"`
}
