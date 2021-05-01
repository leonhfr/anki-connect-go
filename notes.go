package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type (
	// NoteInfo represents information about a single note.
	NoteInfo struct {
		ID     int        `json:"noteId"`    // Unique identifier.
		Model  string     `json:"modelName"` // Model name.
		Tags   []string   `json:"tags"`      // List of tags.
		Fields FieldsInfo `json:"fields"`    // Map of the fields.
	}

	// FieldInfo represents information about a single field.
	FieldInfo struct {
		Value string `json:"value"` // Content of the field.
		Order int    `json:"order"` // Order in which the field appears.
	}

	// FieldsInfo represents information about several fields.
	FieldsInfo struct {
		Front FieldInfo `json:"Front"` // Information about the Front field.
		Back  FieldInfo `json:"Back"`  // Information about the Back field.
	}

	// NoteInput represents a complete note input.
	NoteInput struct {
		Deck    string       `json:"deckName"`  // Name of the deck.
		Model   string       `json:"modelName"` // Name of the model.
		Fields  FieldsInput  `json:"fields"`    // Content of the fields.
		Options OptionsInput `json:"options"`   // Options map.
		Tags    []string     `json:"tags"`      // List of tags.
		Picture []MediaInput `json:"picture"`   // Optional list of picture files.
		Audio   []MediaInput `json:"audio"`     // Optional list of audio files.
		Video   []MediaInput `json:"video"`     // Optional list of video files.
	}

	// FieldsInput represents the fields input.
	FieldsInput struct {
		Front string `json:"Front"` // Update to the Front field.
		Back  string `json:"Back"`  // Update to the Back field.
	}

	// OptionsInput represents the options input.
	OptionsInput struct {
		AllowDuplicate bool `json:"allowDuplicate"` // If true, allows adding duplicate cards. Normally duplicate cards cannot be added and trigger an exception.
	}

	// NoteFieldsInput represents an update to the note's fields input.
	NoteFieldsInput struct {
		ID      int          `json:"id"`      // Note identifier.
		Fields  FieldsInput  `json:"fields"`  // Update to the fields.
		Picture []MediaInput `json:"picture"` // Optional list of picture files.
		Audio   []MediaInput `json:"audio"`   // Optional list of audio files.
		Video   []MediaInput `json:"video"`   // Optional list of video files.
	}

	// MediaInput represents either a picture, an audio, or a video input.
	MediaInput struct {
		URL      string `json:"url"`      // Url to download the file from.
		Filename string `json:"filename"` // Filename to save the media under.
		// TODO: Research skip hash conditions
		SkipHash string   `json:"skipHash"` // Optional md5 to skip inclusion.
		Fields   []string `json:"fields"`   // Optional list of fields to display the media in. Usually Front or Back.
	}
)

// FindNotes returns an array of note ids for a given query.
// Query syntax documentation: https://docs.ankiweb.net/#/searching
func (c *Client) FindNotes(ctx context.Context, query string) ([]int, error) {
	request := findNotesRequest{
		ankiRequest: ankiRequest{Action: "findNotes", Version: c.minVersion},
		Params:      findNotesParams{Query: query},
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

	var res []int
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// NotesInfo returns a list of objects containing information for each given note id.
func (c *Client) NotesInfo(ctx context.Context, notes []int) ([]NoteInfo, error) {
	request := notesInfoRequest{
		ankiRequest: ankiRequest{Action: "notesInfo", Version: c.minVersion},
		Params:      notesInfoParams{Notes: notes},
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

	var res []NoteInfo
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// AddNote creates a note.
// Returns an identifier for the created note, or null if the note couldn't be created.
// Optional picture, audio, or video can be included.
func (c *Client) AddNote(ctx context.Context, note NoteInput) (int, error) {
	request := addNoteRequest{
		ankiRequest: ankiRequest{Action: "addNote", Version: c.minVersion},
		Params:      addNoteParams{Note: note},
	}
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

// AddNotes creates multiple notes.
// Returns an array of ids of the created notes, or null for notes that couldn't be created.
// Optional picture, audio, or video can be included.
func (c *Client) AddNotes(ctx context.Context, notes []NoteInput) ([]int, error) {
	// TODO: implement null in response array: [1496198395707, null]

	request := addNotesRequest{
		ankiRequest: ankiRequest{Action: "addNotes", Version: c.minVersion},
		Params:      addNotesParams{Notes: notes},
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

	var res []int
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateNote modifies the fields of an existing note.
// Optional picture, audio, or video can be included.
func (c *Client) UpdateNote(ctx context.Context, note NoteFieldsInput) error {
	request := updateNoteRequest{
		ankiRequest: ankiRequest{Action: "updateNoteFields", Version: c.minVersion},
		Params:      updateNoteParams{Note: note},
	}
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

// DeleteNotes deletes the notes with the given ids.
// If a note has cards associated with it, all of them will be deleted.
func (c *Client) DeleteNotes(ctx context.Context, notes []int) error {
	request := deleteNotesRequest{
		ankiRequest: ankiRequest{Action: "deleteNotes", Version: c.minVersion},
		Params:      deleteNotesParams{Notes: notes},
	}
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

// DeleteEmptyNotes deletes all the empty notes for the current user.
func (c *Client) DeleteEmptyNotes(ctx context.Context) error {
	request := ankiRequest{Action: "removeEmptyNotes", Version: c.minVersion}
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

type findNotesParams struct {
	Query string `json:"query"`
}

type findNotesRequest struct {
	ankiRequest
	Params struct {
		Query string `json:"query"`
	} `json:"params"`
}

type notesInfoParams struct {
	Notes []int `json:"notes"`
}

type notesInfoRequest struct {
	ankiRequest
	Params notesInfoParams `json:"params"`
}

type addNoteParams struct {
	Note NoteInput `json:"note"`
}

type addNoteRequest struct {
	ankiRequest
	Params addNoteParams `json:"params"`
}

type addNotesParams struct {
	Notes []NoteInput `json:"notes"`
}

type addNotesRequest struct {
	ankiRequest
	Params addNotesParams `json:"params"`
}

type updateNoteParams struct {
	Note NoteFieldsInput `json:"note"`
}

type updateNoteRequest struct {
	ankiRequest
	Params updateNoteParams `json:"params"`
}

type deleteNotesParams struct {
	Notes []int `json:"notes"`
}

type deleteNotesRequest struct {
	ankiRequest
	Params deleteNotesParams `json:"params"`
}
