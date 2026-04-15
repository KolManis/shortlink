package handlers

import (
	"io"
	"net/http"
	"strings"

	urlUseCase "github.com/KolManis/shortlink/internal/usecase/url"
)

type UrlHandler struct {
	usecase urlUseCase.Usecase
}

func NewUrlHandler(usecase urlUseCase.Usecase) *UrlHandler {
	return &UrlHandler{usecase: usecase}
}

// CreateShortURL обработчик POST / - создание короткой ссылки
func (h *UrlHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	// Получаем оригинальный URL
	originalURL := strings.TrimSpace(string(body))
	if originalURL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	// Создаём короткую ссылку через usecase
	shortURL, err := h.usecase.CreateShortURL(r.Context(), originalURL)
	if err != nil {
		http.Error(w, "Failed to create short link", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// Redirect обработчик GET /{id} - редирект на оригинальный URL
func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusBadRequest)
		return
	}

	// Получаем ID из пути
	var shortID string
	if r.URL.Path != "" {
		shortID = strings.TrimPrefix(r.URL.Path, "/")
	}

	if shortID == "" {
		http.Error(w, "Missing short ID", http.StatusBadRequest)
		return
	}

	originalURL, err := h.usecase.GetOriginalURL(r.Context(), shortID)
	if err != nil {
		http.Error(w, "Short URL not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect) // 307
}
