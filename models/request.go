package models

type SnippetChangedRequest struct {
	Snippet Snippet        `json:"snippet"`
	Changes SnippetChanges `json:"changes"`
}
