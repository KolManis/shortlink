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
	// Валидация
	if originalURL == "" {
		return "", ErrInvalidURL
	}

	// Сначала создаём запись без short_code (пока не знаем ID)
	url := &urlDomain.Url{
		OriginalURL: originalURL,
		CreatedAt:   u.now(),
		Clicks:      0,
	}

	// Сохраняем в БД, получаем ID
	dbID, err := u.repo.Create(ctx, url)
	if err != nil {
		return "", fmt.Errorf("failed to create link: %w", err)
	}

	// Генерируем короткий код из ID
	shortCode := encodeBase62(dbID)

	// Обновляем запись с short_code
	url.ID = dbID
	url.ShortCode = shortCode
	if err := u.repo.UpdateShortCode(ctx, dbID, shortCode); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:8080/%s", shortCode), nil
}

func (u *Service) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	link, err := u.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	go func() {
		_ = u.repo.IncrementClicks(context.Background(), shortCode)
	}()

	return link.OriginalURL, nil
}

// encodeBase62 конвертирует число в строку
func encodeBase62(num int64) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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
