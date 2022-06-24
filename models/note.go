package models

type Note struct {
	ID   uint   `json:"-"`
	Text string `json:"text"`
	HTML string `json:"html"`
	Tags []Tag  `json:"tags"`
}
