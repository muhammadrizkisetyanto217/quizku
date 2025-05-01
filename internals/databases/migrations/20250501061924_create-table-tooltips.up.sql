CREATE TABLE IF NOT EXISTS tooltips (
    id SERIAL PRIMARY KEY,
    keyword TEXT NOT NULL UNIQUE CHECK (char_length(keyword) <= 100),
    description_short TEXT NOT NULL CHECK (char_length(description_short) <= 200),
    description_long TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index tambahan untuk pencarian keyword cepat (opsional jika query LIKE digunakan)
CREATE INDEX IF NOT EXISTS idx_tooltips_keyword_lower ON tooltips(LOWER(keyword));
