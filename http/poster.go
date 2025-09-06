package http

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type Poster struct {
	url    string
	bearer string
	client *http.Client
	data   []byte
}

func NewPoster(url string, bearer string, data []byte) *Poster {
	if err := validateArgs(url, bearer); err != nil {
		log.Fatalf("Error validating arguments: %v", err)
	}
	return &Poster{
		url:    url,
		bearer: bearer,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		data:   data,
	}
}

func (p *Poster) Post() (*http.Response, error) {
	req, err := http.NewRequest("POST", p.url, bytes.NewBuffer(p.data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.bearer)

	return p.client.Do(req)
}
