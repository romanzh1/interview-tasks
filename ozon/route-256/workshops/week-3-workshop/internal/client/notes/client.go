package notes

import (
	desc "gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/pkg/api/notes/v1"
)

type Client struct {
	header string
	client desc.NotesClient
}

func NewClient(header string, client desc.NotesClient) *Client {
	return &Client{client: client}
}
