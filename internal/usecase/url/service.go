package url

import (
	"context"
	"fmt"
	"time"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (u *Service) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	// Генерируем короткий ID
	shortID := generateShortID()

	// Создаём модель
	url := &urlDomain.Url{
		ID:          shortID,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		Clicks:      0,
	}

	// Сохраняем в БД
	if err := u.repo.Create(ctx, url); err != nil {
		return "", fmt.Errorf("failed to save link: %w", err)
	}

	// Формируем полную короткую ссылку
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortID)

	return shortURL, nil
}

func (u *Service) GetOriginalURL(ctx context.Context, shortID string) (string, error) {
	// Получаем из БД
	link, err := u.repo.GetByID(ctx, shortID)
	if err != nil {
		return "", fmt.Errorf("failed to get link: %w", err)
	}

	// Увеличиваем счётчик кликов (асинхронно, чтобы не тормозить редирект)
	go func() {
		_ = u.repo.IncrementClicks(context.Background(), shortID)
	}()

	return link.OriginalURL, nil
}

// generateShortID генерирует короткий ID (base62)
func generateShortID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Для простоты используем timestamp
	// В реальном проекте лучше использовать последовательный ID из БД
	timestamp := time.Now().UnixNano()
	num := timestamp % 1000000

	if num == 0 {
		return string(alphabet[0])
	}

	result := make([]byte, 0)
	for num > 0 {
		result = append([]byte{alphabet[num%62]}, result...)
		num /= 62
	}

	return string(result)
}
