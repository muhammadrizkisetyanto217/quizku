CREATE TABLE IF NOT EXISTS categories (
    category_id SERIAL PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL,
    category_status VARCHAR(10) CHECK (category_status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    category_description_short VARCHAR(100),
    category_description_long VARCHAR(2000),
    category_total_subcategories INTEGER[] NOT NULL DEFAULT '{}',
    category_image_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    category_difficulty_id INT REFERENCES difficulties(difficulty_id),

    CONSTRAINT unique_category_name UNIQUE (category_name)
);

-- ✅ Indexing
CREATE INDEX IF NOT EXISTS idx_category_difficulty_id ON categories(category_difficulty_id);
CREATE INDEX IF NOT EXISTS idx_category_status ON categories(category_status);



CREATE TABLE IF NOT EXISTS categories_news (
    category_news_id SERIAL PRIMARY KEY,
    category_news_category_id INTEGER NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
    category_news_title VARCHAR(255) NOT NULL,
    category_news_description TEXT NOT NULL,
    category_news_is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Indexing
CREATE INDEX IF NOT EXISTS idx_category_news_category_id ON categories_news(category_news_category_id);
CREATE INDEX IF NOT EXISTS idx_category_news_is_public ON categories_news(category_news_is_public);
CREATE INDEX IF NOT EXISTS idx_news_public_per_category ON categories_news(category_news_category_id, category_news_is_public);
