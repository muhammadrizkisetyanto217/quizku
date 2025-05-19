-- ✅ TABLE: units
CREATE TABLE IF NOT EXISTS units (
    unit_id SERIAL PRIMARY KEY,
    unit_name VARCHAR(50) UNIQUE NOT NULL,
    unit_status VARCHAR(10) CHECK (unit_status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    unit_description_short VARCHAR(200) NOT NULL,
    unit_description_overview TEXT NOT NULL,
    unit_image_url VARCHAR(100),
    unit_total_section_quizzes INTEGER[] NOT NULL DEFAULT '{}',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    unit_themes_or_level_id INT REFERENCES themes_or_levels(themes_or_level_id) ON DELETE CASCADE,
    unit_created_by UUID REFERENCES users(user_id) ON DELETE CASCADE
);

-- ✅ Indexing for performance
CREATE INDEX IF NOT EXISTS idx_unit_status ON units(unit_status);
CREATE INDEX IF NOT EXISTS idx_unit_themes_id ON units(unit_themes_or_level_id);
CREATE INDEX IF NOT EXISTS idx_unit_created_by ON units(unit_created_by);


-- ✅ TABLE: units_news
-- ✅ TABLE: unit_news
CREATE TABLE IF NOT EXISTS unit_news (
    unit_news_id SERIAL PRIMARY KEY,
    unit_news_unit_id INTEGER NOT NULL REFERENCES units(unit_id) ON DELETE CASCADE,
    unit_news_title VARCHAR(255) NOT NULL,
    unit_news_description TEXT NOT NULL,
    unit_news_is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Indexing for units_news
CREATE INDEX IF NOT EXISTS idx_unit_news_unit_id ON unit_news(unit_news_unit_id);
CREATE INDEX IF NOT EXISTS idx_unit_news_is_public ON unit_news(unit_news_is_public);
CREATE INDEX IF NOT EXISTS idx_unit_news_unit_public ON unit_news(unit_news_unit_id, unit_news_is_public);



CREATE TABLE IF NOT EXISTS user_unit (
    user_unit_id SERIAL PRIMARY KEY,
    user_unit_user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    user_unit_unit_id INTEGER NOT NULL REFERENCES units(unit_id) ON DELETE CASCADE,

    user_unit_attempt_reading INTEGER NOT NULL DEFAULT 0,
    user_unit_attempt_evaluation JSONB NOT NULL DEFAULT '{}'::jsonb,
    user_unit_complete_section_quizzes JSONB NOT NULL DEFAULT '{}'::jsonb,

    user_unit_grade_quiz INTEGER NOT NULL DEFAULT 0 CHECK (user_unit_grade_quiz BETWEEN 0 AND 100),
    user_unit_grade_exam INTEGER NOT NULL DEFAULT 0 CHECK (user_unit_grade_exam BETWEEN 0 AND 100),
    user_unit_grade_result INTEGER NOT NULL DEFAULT 0 CHECK (user_unit_grade_result BETWEEN 0 AND 100),
    user_unit_is_passed BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Indexing untuk performa
CREATE INDEX IF NOT EXISTS idx_user_unit_user_unit_user_id_unit_id ON user_unit (user_unit_user_id, user_unit_unit_id);
CREATE INDEX IF NOT EXISTS idx_user_unit_user_id ON user_unit (user_unit_user_id);
CREATE INDEX IF NOT EXISTS idx_user_unit_unit_id ON user_unit (user_unit_unit_id);
