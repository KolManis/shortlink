package url

import "time"

type Url struct {
	ID          string
	OriginalURL string
	CreatedAt   time.Time
	Clicks      int64
}
