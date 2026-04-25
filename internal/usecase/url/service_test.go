package url

import (
	"context"
	"log/slog"
	"testing"
	"time"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetOriginalURL_Success(t *testing.T) {
	mockRepo := NewMockRepository(t)
	mockCache := NewMockCache(t)

	service := &Service{
		repo:   mockRepo,
		cache:  mockCache,
		logger: slog.Default(),
	}

	ctx := context.Background()
	shortCode := "abc123"
	cacheKey := "url:" + shortCode
	expectedURL := "https://example.com"

	mockCache.EXPECT().
		Get(mock.Anything, cacheKey).
		Return("", urlDomain.ErrNotFound). // кэш пуст
		Once()

	mockRepo.EXPECT().
		GetByShortCode(mock.Anything, shortCode).
		Return(&urlDomain.Url{
			ShortCode:   shortCode,
			OriginalURL: expectedURL,
		}, nil).
		Once()

	mockCache.EXPECT().
		Set(mock.Anything, cacheKey, expectedURL, time.Hour).
		Return(nil).
		Once()

	result, err := service.GetOriginalURL(ctx, shortCode)

	require.NoError(t, err)
	assert.Equal(t, expectedURL, result)
}

func TestService_GetOriginalURL_CacheHit(t *testing.T) {
	mockRepo := NewMockRepository(t)
	mockCache := NewMockCache(t)

	service := &Service{
		repo:   mockRepo,
		cache:  mockCache,
		logger: slog.Default(),
	}

	ctx := context.Background()
	shortCode := "abc123"
	cacheKey := "url:" + shortCode
	expectedURL := "https://example.com"

	// Данные есть в кэше
	mockCache.EXPECT().
		Get(mock.Anything, cacheKey).
		Return(expectedURL, nil). // ← успешно вернули URL
		Once()

	// Репозиторий НЕ ДОЛЖЕН вызываться!
	// mockRepo.EXPECT().GetByShortCode(...).Times(0)

	result, err := service.GetOriginalURL(ctx, shortCode)

	require.NoError(t, err)
	assert.Equal(t, expectedURL, result)
}
