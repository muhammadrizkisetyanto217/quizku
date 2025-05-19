CREATE TABLE IF NOT EXISTS exams (
    exam_id SERIAL PRIMARY KEY,
    exam_name VARCHAR(100) NOT NULL,
    exam_status VARCHAR(20) CHECK (exam_status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    exam_total_question_ids INTEGER[] NOT NULL DEFAULT '{}',
    exam_icon_url VARCHAR(255),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP

    exam_unit_id INT REFERENCES units(unit_id) ON DELETE CASCADE,
    exam_created_by UUID REFERENCES users(user_id) ON DELETE CASCADE
);

-- ✅ Index untuk performa pencarian
CREATE INDEX IF NOT EXISTS idx_exam_status ON exams(exam_status);
CREATE INDEX IF NOT EXISTS idx_exam_unit_id ON exams(exam_unit_id);
CREATE INDEX IF NOT EXISTS idx_exam_created_by ON exams(exam_created_by);



CREATE TABLE IF NOT EXISTS user_exams (
    user_exam_id SERIAL PRIMARY KEY,

    user_exam_user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    user_exam_exam_id INTEGER NOT NULL REFERENCES exams(exam_id) ON DELETE CASCADE,
    user_exam_unit_id INTEGER NOT NULL REFERENCES units(unit_id) ON DELETE CASCADE,

    user_exam_attempt INTEGER NOT NULL DEFAULT 1,
    user_exam_percentage_grade INTEGER NOT NULL DEFAULT 0,
    user_exam_time_duration INTEGER NOT NULL DEFAULT 0,
    user_exam_point INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Index dengan nama dan kolom yang semantik
CREATE INDEX IF NOT EXISTS idx_user_exam_user_exam_id ON user_exams (user_exam_user_id, user_exam_exam_id);
CREATE INDEX IF NOT EXISTS idx_user_exam_user_unit_id ON user_exams (user_exam_user_id, user_exam_unit_id);
