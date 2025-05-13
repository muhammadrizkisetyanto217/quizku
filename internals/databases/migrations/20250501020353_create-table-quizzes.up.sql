-- ✅ TABLE: section_quizzes
CREATE TABLE IF NOT EXISTS section_quizzes (
    id SERIAL PRIMARY KEY,
    name_section_quizzes VARCHAR(50) NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    materials_quizzes TEXT NOT NULL,
    icon_url VARCHAR(100),
    total_quizzes INTEGER[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    unit_id INT REFERENCES units(id) ON DELETE CASCADE,
    created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- ✅ Indexing section_quizzes
CREATE INDEX IF NOT EXISTS idx_section_quizzes_status ON section_quizzes(status);
CREATE INDEX IF NOT EXISTS idx_section_quizzes_unit_id ON section_quizzes(unit_id);
CREATE INDEX IF NOT EXISTS idx_section_quizzes_created_by ON section_quizzes(created_by);
CREATE INDEX IF NOT EXISTS idx_section_unit_status ON section_quizzes(unit_id, status);



CREATE TABLE IF NOT EXISTS user_section_quizzes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    section_quizzes_id INTEGER NOT NULL,
    complete_quiz JSONB NOT NULL DEFAULT '{}'::jsonb,
    grade_result INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_section_quizzes_user_id ON user_section_quizzes (user_id);
CREATE INDEX IF NOT EXISTS idx_user_section_quizzes_section_id ON user_section_quizzes (section_quizzes_id);


-- ✅ TABLE: quizzes
CREATE TABLE IF NOT EXISTS quizzes (
    id SERIAL PRIMARY KEY,
    name_quizzes VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    total_question INTEGER[] NOT NULL DEFAULT '{}',
    icon_url VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    section_quizzes_id INT REFERENCES section_quizzes(id) ON DELETE CASCADE,
    created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- ✅ Indexing quizzes
CREATE INDEX IF NOT EXISTS idx_quizzes_status ON quizzes(status);
CREATE INDEX IF NOT EXISTS idx_quizzes_section_id ON quizzes(section_quizzes_id);
CREATE INDEX IF NOT EXISTS idx_quizzes_created_by ON quizzes(created_by);



CREATE TABLE IF NOT EXISTS user_quizzes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    quiz_id INTEGER NOT NULL,
    attempt INTEGER NOT NULL DEFAULT 1,
    percentage_grade INTEGER NOT NULL DEFAULT 0,
    time_duration INTEGER NOT NULL DEFAULT 0,
    point INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_quizzes_user_id ON user_quizzes (user_id);
CREATE INDEX IF NOT EXISTS idx_user_quizzes_quiz_id ON user_quizzes (quiz_id);