-- ✅ TABLE: difficulties
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

-- ✅ Indexing
CREATE INDEX IF NOT EXISTS idx_difficulty_status ON difficulties(difficulty_status);

-- ✅ TABLE: difficulties_news
CREATE TABLE IF NOT EXISTS difficulties_news (
    difficulty_news_id SERIAL PRIMARY KEY,
    difficulty_news_difficulty_id INTEGER NOT NULL REFERENCES difficulties(difficulty_id) ON DELETE CASCADE,
    difficulty_news_title VARCHAR(255) NOT NULL,
    difficulty_news_description TEXT NOT NULL,
    difficulty_news_is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Indexes
CREATE INDEX IF NOT EXISTS idx_difficulty_news_difficulty_id ON difficulties_news(difficulty_news_difficulty_id);
CREATE INDEX IF NOT EXISTS idx_difficulty_news_is_public ON difficulties_news(difficulty_news_is_public);
CREATE INDEX IF NOT EXISTS idx_difficulty_news_combined ON difficulties_news(difficulty_news_difficulty_id, difficulty_news_is_public);
