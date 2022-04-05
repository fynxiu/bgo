package gohttp

import (
	"io"
	"net/http"

)


func Do(req *http.Request) (resp *http.Response, err error) {
	return http.DefaultClient.Do(req)
}

// Get sends an HTTP GET request to the service.
func (s namedService) Get(url string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return s.Do(r)
}

// Post sends an HTTP POST request to the service.
func Post(url, ctype string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set(contentType, ctype)
	return s.Do(r)
}


		resp, err = s.cli.Do(r)
		return err

	return
}
