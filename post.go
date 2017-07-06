package main

type Post struct {
	Id  int    `json:"id"`
	URI string `json:"uri"`

	Title string   `json:"title"`
	Size  int      `json:"size"`
	Files []string `json:"files"`

	Time int `json:"time"`

	Tags []string          `json:"tags"`
	Meta map[string]string `json:"meta"`
}
