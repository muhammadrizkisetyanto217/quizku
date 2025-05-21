-- ✅ TABLE: section_quizzes (Refactored)
CREATE TABLE IF NOT EXISTS section_quizzes (
    section_quizzes_id SERIAL PRIMARY KEY,
    section_quizzes_name VARCHAR(50) NOT NULL,
    section_quizzes_status VARCHAR(10) CHECK (section_quizzes_status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    section_quizzes_materials TEXT NOT NULL,
    section_quizzes_icon_url VARCHAR(100),
    section_quizzes_total_quizzes INTEGER[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    section_quizzes_unit_id INT REFERENCES units(unit_id) ON DELETE CASCADE,
    section_quizzes_created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- ✅ Indexing dengan nama yang deskriptif
CREATE INDEX IF NOT EXISTS idx_section_quizzes_status ON section_quizzes(section_quizzes_status);
CREATE INDEX IF NOT EXISTS idx_section_quizzes_unit_id ON section_quizzes(section_quizzes_unit_id);
CREATE INDEX IF NOT EXISTS idx_section_quizzes_created_by ON section_quizzes(section_quizzes_created_by);
CREATE INDEX IF NOT EXISTS idx_section_quizzes_unit_status ON section_quizzes(section_quizzes_unit_id, section_quizzes_status);


-- ✅ REFACTORED: Struktur tabel user_section_quizzes dengan nama kolom deskriptif
CREATE TABLE IF NOT EXISTS user_section_quizzes (
    user_section_quizzes_id SERIAL PRIMARY KEY,
    user_section_quizzes_user_id UUID NOT NULL,
    user_section_quizzes_section_quizzes_id INTEGER NOT NULL,
    user_section_quizzes_complete_quiz JSONB NOT NULL DEFAULT '{}'::jsonb,
    user_section_quizzes_grade_result INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Index untuk performa query berdasarkan user dan section
CREATE INDEX IF NOT EXISTS idx_user_section_quizzes_user_id 
    ON user_section_quizzes (user_section_quizzes_user_id);

CREATE INDEX IF NOT EXISTS idx_user_section_quizzes_section_id 
    ON user_section_quizzes (user_section_quizzes_section_quizzes_id);



-- ✅ TABLE: quizzes
CREATE TABLE IF NOT EXISTS quizzes (
    quiz_id SERIAL PRIMARY KEY,
    quiz_name VARCHAR(50) UNIQUE NOT NULL,
    quiz_status VARCHAR(10) CHECK (quiz_status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    quiz_total_question INTEGER[] NOT NULL DEFAULT '{}',
    quiz_icon_url VARCHAR(100),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    quiz_section_quizzes_id INT REFERENCES section_quizzes(section_quizzes_id) ON DELETE CASCADE,
    quiz_created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- ✅ Indexes
CREATE INDEX IF NOT EXISTS idx_quizzes_status ON quizzes(quiz_status);
CREATE INDEX IF NOT EXISTS idx_quizzes_section_quizzes_id ON quizzes(quiz_section_quizzes_id);
CREATE INDEX IF NOT EXISTS idx_quizzes_created_by ON quizzes(quiz_created_by);


CREATE TABLE IF NOT EXISTS user_quizzes (
    user_quiz_id SERIAL PRIMARY KEY,
    user_quiz_user_id UUID NOT NULL,
    user_quiz_quiz_id INTEGER NOT NULL,
    user_quiz_attempt INTEGER NOT NULL DEFAULT 1,
    user_quiz_percentage_grade INTEGER NOT NULL DEFAULT 0,
    user_quiz_time_duration INTEGER NOT NULL DEFAULT 0,
    user_quiz_point INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Indexing
CREATE INDEX IF NOT EXISTS idx_user_quizzes_user_id ON user_quizzes (user_quiz_user_id);
CREATE INDEX IF NOT EXISTS idx_user_quizzes_quiz_id ON user_quizzes (user_quiz_quiz_id);