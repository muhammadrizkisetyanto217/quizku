-- ✅ TABLE: themes_or_levels
-- ✅ TABLE: themes_or_levels (semantik)
CREATE TABLE IF NOT EXISTS themes_or_levels (
    themes_or_level_id SERIAL PRIMARY KEY,
    themes_or_level_name VARCHAR(255) NOT NULL,
    themes_or_level_status VARCHAR(10) CHECK (themes_or_level_status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    themes_or_level_description_short VARCHAR(100),
    themes_or_level_description_long VARCHAR(2000),
    themes_or_level_total_unit INTEGER[] NOT NULL DEFAULT '{}',
    themes_or_level_image_url VARCHAR(100),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    themes_or_level_subcategory_id INT REFERENCES subcategories(id) ON DELETE SET NULL,

    CONSTRAINT unique_theme_name_per_subcategory UNIQUE (themes_or_level_name, themes_or_level_subcategory_id)
);

-- ✅ Indexing themes_or_levels (semantik)
CREATE INDEX IF NOT EXISTS idx_themes_or_level_status ON themes_or_levels(themes_or_level_status);
CREATE INDEX IF NOT EXISTS idx_themes_or_level_subcategory_id ON themes_or_levels(themes_or_level_subcategory_id);
CREATE INDEX IF NOT EXISTS idx_themes_or_level_name_subcat ON themes_or_levels(themes_or_level_name, themes_or_level_subcategory_id);


-- ✅ TABLE: themes_or_levels_news
CREATE TABLE IF NOT EXISTS themes_or_levels_news (
    themes_news_id SERIAL PRIMARY KEY,
    themes_news_themes_or_level_id INTEGER NOT NULL REFERENCES themes_or_levels(themes_or_level_id) ON DELETE CASCADE,
    themes_news_title VARCHAR(255) NOT NULL,
    themes_news_description TEXT NOT NULL,
    themes_news_is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Indexing themes_or_levels_news
CREATE INDEX IF NOT EXISTS idx_themes_news_theme_id ON themes_or_levels_news(themes_news_themes_or_level_id);
CREATE INDEX IF NOT EXISTS idx_themes_news_is_public ON themes_or_levels_news(themes_news_is_public);
CREATE INDEX IF NOT EXISTS idx_themes_news_per_theme_public 
  ON themes_or_levels_news(themes_news_themes_or_level_id, themes_news_is_public);



-- ✅ TABLE: user_themes_or_levels
CREATE TABLE IF NOT EXISTS user_themes_or_levels (
    user_theme_id SERIAL PRIMARY KEY,
    user_theme_user_id UUID NOT NULL,
    user_theme_themes_or_level_id INTEGER NOT NULL REFERENCES themes_or_levels(themes_or_level_id) ON DELETE CASCADE,
    user_theme_complete_unit JSONB NOT NULL DEFAULT '{}'::jsonb,
    user_theme_grade_result INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Indexing user_themes_or_levels
CREATE INDEX IF NOT EXISTS idx_user_theme_user_id ON user_themes_or_levels (user_theme_user_id);
CREATE INDEX IF NOT EXISTS idx_user_theme_theme_id ON user_themes_or_levels (user_theme_themes_or_level_id);
CREATE INDEX IF NOT EXISTS idx_user_theme_unique ON user_themes_or_levels(user_theme_user_id, user_theme_themes_or_level_id);
