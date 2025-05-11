-- ✅ TABLE: themes_or_levels
CREATE TABLE IF NOT EXISTS themes_or_levels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    description_short VARCHAR(100),
    description_long VARCHAR(2000),
    total_unit INTEGER[] NOT NULL DEFAULT '{}',
    image_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    subcategories_id INT REFERENCES subcategories(id) ON DELETE SET NULL,

    CONSTRAINT unique_theme_name_per_subcategory UNIQUE (name, subcategories_id)
);

-- ✅ Indexing themes_or_levels
CREATE INDEX IF NOT EXISTS idx_themes_status ON themes_or_levels(status);
CREATE INDEX IF NOT EXISTS idx_themes_subcategories_id ON themes_or_levels(subcategories_id);
CREATE INDEX IF NOT EXISTS idx_themes_name_subcat ON themes_or_levels(name, subcategories_id);

-- ✅ TABLE: themes_or_levels_news
CREATE TABLE IF NOT EXISTS themes_or_levels_news (
    id SERIAL PRIMARY KEY,
    themes_or_levels_id INTEGER NOT NULL REFERENCES themes_or_levels(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Indexing themes_or_levels_news
CREATE INDEX IF NOT EXISTS idx_themes_news_theme_id ON themes_or_levels_news(themes_or_levels_id);
CREATE INDEX IF NOT EXISTS idx_themes_news_is_public ON themes_or_levels_news(is_public);
CREATE INDEX IF NOT EXISTS idx_themes_news_per_theme_public ON themes_or_levels_news(themes_or_levels_id, is_public);

-- ✅ TABLE: user_themes_or_levels
CREATE TABLE IF NOT EXISTS user_themes_or_levels (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    themes_or_levels_id INTEGER NOT NULL REFERENCES themes_or_levels(id) ON DELETE CASCADE,
    complete_unit JSONB NOT NULL DEFAULT '{}'::jsonb,
    grade_result INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Indexing user_themes_or_levels
CREATE INDEX IF NOT EXISTS idx_user_themes_user_id ON user_themes_or_levels (user_id);
CREATE INDEX IF NOT EXISTS idx_user_themes_theme_id ON user_themes_or_levels (themes_or_levels_id);
CREATE INDEX IF NOT EXISTS idx_user_themes_user_theme ON user_themes_or_levels(user_id, themes_or_levels_id);
