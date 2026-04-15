package url

import "time"

type Url struct {
	ID          int64
	ShortCode   string
	OriginalURL string
	CreatedAt   time.Time
	Clicks      int64
}
