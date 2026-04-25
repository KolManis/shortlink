package url

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
	"github.com/KolManis/shortlink/internal/repository/postgres"
	"github.com/KolManis/shortlink/internal/uow"
)

type Service struct {
	repo   Repository
	uow    *uow.UnitOfWork
	cache  Cache
	logger *slog.Logger
	now    func() time.Time
}

func NewService(repo Repository, uow *uow.UnitOfWork, cache Cache, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		uow:    uow,
		cache:  cache,
		logger: logger,
		now:    func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	s.logger.Debug("creating short URL", "original_url", originalURL)

	if originalURL == "" {
		s.logger.Warn("empty URL provided")
		return "", ErrInvalidURL
	}

	existing, err := s.repo.GetByOriginalURL(ctx, originalURL)

	if err == nil && existing != nil {
		s.logger.Info("URL already exists", "original_url", originalURL, "short_code", existing.ShortCode)
		return fmt.Sprintf("http://localhost:8080/%s", existing.ShortCode), nil
	}

	if err != nil && !errors.Is(err, urlDomain.ErrNotFound) {
		s.logger.Error("failed to check existing URL", "error", err)
		return "", fmt.Errorf("failed to check existing: %w", err)
	}

	var shortCode string

	err = s.uow.Do(ctx, func(r *postgres.Repository) error {
		url := &urlDomain.Url{
			OriginalURL: originalURL,
			CreatedAt:   s.now(),
			Clicks:      0,
		}

		id, err := r.Create(ctx, url)
		if err != nil {
			return err
		}

		shortCode = encodeBase62(id)

		return r.UpdateShortCode(ctx, id, shortCode)
	})

	if err != nil {
		s.logger.Error("failed to create short URL", "error", err)
		return "", err
	}

	s.logger.Info("short URL created", "short_code", shortCode)
	return fmt.Sprintf("http://localhost:8080/%s", shortCode), nil
}

func (s *Service) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	s.logger.Debug("getting original URL", "short_code", shortCode)

	cached, err := s.cache.Get(ctx, "url:"+shortCode)
	if err == nil {
		s.logger.Debug("cache hit", "short_code", shortCode, "url", cached)
		return cached, nil
	}

	s.logger.Debug("cache miss", "short_code", shortCode)

	link, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		s.logger.Warn("short code not found", "short_code", shortCode, "error", err)
		return "", err
	}

	if err := s.cache.Set(ctx, "url:"+shortCode, link.OriginalURL, time.Hour); err != nil {
		s.logger.Warn("failed to set cache", "short_code", shortCode, "error", err)
	}

	// go func() {
	// 	_ = u.repo.IncrementClicks(context.Background(), shortCode)
	// }()

	// s.logger.Info("redirect", "short_code", shortCode, "original_url", link.OriginalURL)

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
