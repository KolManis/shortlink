package transporthttp

import (
	"net/http"

	httpHandlers "github.com/KolManis/shortlink/internal/transport/http/handlers"
)

func NewRouter(urlHandler *httpHandlers.UrlHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", urlHandler.CreateShortURL)
	mux.HandleFunc("POST /api/shorten", urlHandler.CreateShortURLJSON)
	mux.HandleFunc("GET /{id}", urlHandler.Redirect)

	return mux
}
