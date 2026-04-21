package url

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
)

type Service struct {
	repo   Repository
	cache  Cache
	logger *slog.Logger
	now    func() time.Time
}

func NewService(repo Repository, cache Cache, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		cache:  cache,
		logger: logger,
		now:    func() time.Time { return time.Now().UTC() },
	}
}

func (u *Service) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	u.logger.Debug("creating short URL", "original_url", originalURL)

	if originalURL == "" {
		u.logger.Warn("empty URL provided")
		return "", ErrInvalidURL
	}

	existing, err := u.repo.GetByOriginalURL(ctx, originalURL)

	if err == nil && existing != nil {
		u.logger.Info("URL already exists", "original_url", originalURL, "short_code", existing.ShortCode)
		return fmt.Sprintf("http://localhost:8080/%s", existing.ShortCode), nil
	}

	if err != nil && !errors.Is(err, urlDomain.ErrNotFound) {
		u.logger.Error("failed to check existing URL", "error", err)
		return "", fmt.Errorf("failed to check existing: %w", err)
	}

	tx, err := u.repo.BeginTx(ctx)
	if err != nil {
		u.logger.Error("failed to begin transaction", "error", err)
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Сначала создаём запись без short_code (пока не знаем ID)
	url := &urlDomain.Url{
		OriginalURL: originalURL,
		CreatedAt:   u.now(),
		Clicks:      0,
	}

	// Сохраняем в БД, получаем ID
	dbID, err := u.repo.Create(ctx, tx, url)
	if err != nil {
		u.logger.Error("failed to create link", "error", err)
		return "", fmt.Errorf("failed to create link: %w", err)
	}

	// Генерируем короткий код из ID
	shortCode := encodeBase62(dbID)

	if err := u.repo.UpdateShortCode(ctx, tx, dbID, shortCode); err != nil {
		u.logger.Error("failed to update short code", "error", err)
		return "", fmt.Errorf("failed to update short code: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		u.logger.Error("failed to commit transaction", "error", err)
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.logger.Info("short URL created", "short_code", shortCode, "original_url", originalURL)

	return fmt.Sprintf("http://localhost:8080/%s", shortCode), nil
}

func (u *Service) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	u.logger.Debug("getting original URL", "short_code", shortCode)

	cached, err := u.cache.Get(ctx, "url:"+shortCode)
	if err == nil {
		u.logger.Debug("cache hit", "short_code", shortCode, "url", cached)
		return cached, nil
	}

	u.logger.Debug("cache miss", "short_code", shortCode)

	link, err := u.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		u.logger.Warn("short code not found or failed to get from DB", "short_code", shortCode)
		return "", err
	}

	if err := u.cache.Set(ctx, "url:"+shortCode, link.OriginalURL, time.Hour); err != nil {
		u.logger.Warn("failed to set cache", "short_code", shortCode, "error", err)
	}

	// go func() {
	// 	_ = u.repo.IncrementClicks(context.Background(), shortCode)
	// }()

	u.logger.Info("redirect", "short_code", shortCode, "original_url", link.OriginalURL)

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
