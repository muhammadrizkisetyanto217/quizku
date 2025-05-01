-- ✅ TABLE: categories (dengan INTEGER ARRAY dan tanpa trigger)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    description_short VARCHAR(100),
    description_long VARCHAR(2000),
    total_subcategories INTEGER[] NOT NULL DEFAULT '{}',
    image_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    difficulty_id INT REFERENCES difficulties(id),

    CONSTRAINT unique_category_name UNIQUE (name)
);

-- ✅ Indexing untuk performa
CREATE INDEX IF NOT EXISTS idx_categories_difficulty_id ON categories(difficulty_id);
CREATE INDEX IF NOT EXISTS idx_categories_status ON categories(status);

-- ✅ TABLE: categories_news
CREATE TABLE IF NOT EXISTS categories_news (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);


-- ✅ Indexing untuk performa
CREATE INDEX IF NOT EXISTS idx_categories_news_category_id ON categories_news(category_id);
CREATE INDEX IF NOT EXISTS idx_categories_news_is_public ON categories_news(is_public);
CREATE INDEX IF NOT EXISTS idx_news_public_per_category ON categories_news(category_id, is_public);

-- ✅ TABLE: user_category
CREATE TABLE IF NOT EXISTS user_category (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    complete_category INTEGER[] NOT NULL DEFAULT '{}',
    total_category INTEGER[] NOT NULL DEFAULT '{}',
    grade_result INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE INDEX IF NOT EXISTS idx_user_category_user_id ON user_category(user_id);
CREATE INDEX IF NOT EXISTS idx_user_category_cat_id ON user_category(category_id);
CREATE INDEX IF NOT EXISTS idx_user_category_user_cat ON user_category(user_id, category_id);
