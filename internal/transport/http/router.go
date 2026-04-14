package transporthttp

import (
	"net/http"
)

func NewRouter(linkHandler *httphandlers.LinkHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", linkHandler.CreateShortURL)
	mux.HandleFunc("GET /{id}", linkHandler.Redirect)

	return mux
}
