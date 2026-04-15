CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,    
    short_code VARCHAR(10) UNIQUE NOT NULL, 
    original_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    clicks BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_short_code ON links(short_code);