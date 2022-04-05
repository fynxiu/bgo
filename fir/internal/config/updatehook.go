package config

import "net/http"

type UpdatedHook struct {
	URL    string      `json:"url" yaml:"url"`
	Header http.Header `json:"header" yaml:"header"`
}
