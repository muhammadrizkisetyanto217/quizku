-- ✅ TABLE: units
CREATE TABLE IF NOT EXISTS units (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    description_short VARCHAR(200) NOT NULL,
    description_overview TEXT NOT NULL,
    image_url VARCHAR(100),
    total_section_quizzes INTEGER[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    themes_or_level_id INT REFERENCES themes_or_levels(id) ON DELETE CASCADE,
    created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- ✅ Indexing for performance
CREATE INDEX IF NOT EXISTS idx_units_status ON units(status);
CREATE INDEX IF NOT EXISTS idx_units_themes_id ON units(themes_or_level_id);
CREATE INDEX IF NOT EXISTS idx_units_created_by ON units(created_by);

-- ✅ TABLE: units_news
CREATE TABLE IF NOT EXISTS units_news (
    id SERIAL PRIMARY KEY,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Indexing for units_news
CREATE INDEX IF NOT EXISTS idx_units_news_unit_id ON units_news(unit_id);
CREATE INDEX IF NOT EXISTS idx_units_news_is_public ON units_news(is_public);
CREATE INDEX IF NOT EXISTS idx_units_news_unit_public ON units_news(unit_id, is_public);
