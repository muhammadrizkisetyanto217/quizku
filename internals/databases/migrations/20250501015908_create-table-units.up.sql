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


CREATE TABLE IF NOT EXISTS user_unit (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    attempt_reading INTEGER DEFAULT 0 NOT NULL,
    attempt_evaluation JSONB NOT NULL DEFAULT '{}'::jsonb,
    complete_section_quizzes JSONB NOT NULL DEFAULT '{}'::jsonb,
    total_section_quizzes INTEGER[] NOT NULL DEFAULT '{}',
    grade_quiz INTEGER NOT NULL DEFAULT 0 CHECK (grade_quiz BETWEEN 0 AND 100),
    grade_exam INTEGER NOT NULL DEFAULT 0 CHECK (grade_exam BETWEEN 0 AND 100),
    grade_result INTEGER NOT NULL DEFAULT 0 CHECK (grade_result BETWEEN 0 AND 100),
    is_passed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_unit_user_id_unit_id ON user_unit (user_id, unit_id);
CREATE INDEX IF NOT EXISTS idx_user_unit_user_id ON user_unit (user_id);
CREATE INDEX IF NOT EXISTS idx_user_unit_unit_id ON user_unit (unit_id);
