package backend

import "net/http"

func NewHTTPClient() *http.Client {
	return http.DefaultClient
}
