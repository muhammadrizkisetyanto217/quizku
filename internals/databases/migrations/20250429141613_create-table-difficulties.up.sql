CREATE TABLE IF NOT EXISTS difficulties (
    difficulty_id SERIAL PRIMARY KEY,
    difficulty_name VARCHAR(255) NOT NULL,
    difficulty_status VARCHAR(10) DEFAULT 'pending' CHECK (difficulty_status IN ('active', 'pending', 'archived')),
    difficulty_description_short VARCHAR(200),
    difficulty_description_long VARCHAR(3000),
    difficulty_total_categories INTEGER[] NOT NULL DEFAULT '{}',
    difficulty_image_url VARCHAR(255),
    difficulty_update_news JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Indexing untuk performa
CREATE INDEX IF NOT EXISTS idx_difficulty_status ON difficulties(difficulty_status);


CREATE TABLE IF NOT EXISTS difficulty_news (
    difficulty_news_id SERIAL PRIMARY KEY,
    difficulty_news_difficulty_id INTEGER NOT NULL REFERENCES difficulties(difficulty_id) ON DELETE CASCADE,
    difficulty_news_title VARCHAR(255) NOT NULL,
    difficulty_news_description TEXT NOT NULL,
    difficulty_news_is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Index untuk pencarian dan filter
CREATE INDEX IF NOT EXISTS idx_difficulty_news_difficulty_id ON difficulty_news(difficulty_news_difficulty_id);
CREATE INDEX IF NOT EXISTS idx_difficulty_news_is_public ON difficulty_news(difficulty_news_is_public);
CREATE INDEX IF NOT EXISTS idx_news_public_per_difficulty ON difficulty_news(difficulty_news_difficulty_id, difficulty_news_is_public);

