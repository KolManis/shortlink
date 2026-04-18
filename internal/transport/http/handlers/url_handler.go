package handlers

import (
	"encoding/json"
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

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
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

func (h *UrlHandler) CreateShortURLJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	// Декодируем JSON
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL field is required", http.StatusBadRequest)
		return
	}

	shortURL, err := h.usecase.CreateShortURL(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "Failed to create short link", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(shortenResponse{Result: shortURL}); err != nil {
		return
	}
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
