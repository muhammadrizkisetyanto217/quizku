-- ✅ TABLE: section_quizzes
CREATE TABLE IF NOT EXISTS section_quizzes (
    id SERIAL PRIMARY KEY,
    name_quizzes VARCHAR(50) NOT NULL,
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
