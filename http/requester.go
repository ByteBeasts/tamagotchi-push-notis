package http

import (
	"log"
	"net/http"
	"time"
)

type Requester struct {
	url    string
	bearer string
	client *http.Client
}

func NewRequester(url string, bearer string) *Requester {
	if err := validateArgs(url, bearer); err != nil {
		log.Fatalf("Error validating arguments: %v", err)
	}
	return &Requester{
		url:    url,
		bearer: bearer,
		client: &http.Client{
			Timeout: 180 * time.Second,
		},
	}
}

func (r *Requester) Request() (*http.Response, error) {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+r.bearer)

	return r.client.Do(req)
}

