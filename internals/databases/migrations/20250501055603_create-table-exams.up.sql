CREATE TABLE IF NOT EXISTS exams (
	id SERIAL PRIMARY KEY,
	name_exams VARCHAR(50) NOT NULL,
	status VARCHAR(10) CHECK (status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
	total_question INTEGER[] NOT NULL DEFAULT '{}',
	icon_url VARCHAR(100),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP,
	unit_id INT REFERENCES units(id) ON DELETE CASCADE,
	created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- Indexing untuk performa
CREATE INDEX IF NOT EXISTS idx_exams_status ON exams(status);
CREATE INDEX IF NOT EXISTS idx_exams_unit_id ON exams(unit_id);
CREATE INDEX IF NOT EXISTS idx_exams_created_by ON exams(created_by);


CREATE TABLE IF NOT EXISTS user_exams (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exam_id INTEGER NOT NULL REFERENCES exams(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    attempt INTEGER NOT NULL DEFAULT 1,
    percentage_grade INTEGER NOT NULL DEFAULT 0,
    time_duration INTEGER NOT NULL DEFAULT 0,
    point INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_exam_user_exam ON user_exams (user_id, exam_id);
CREATE INDEX IF NOT EXISTS idx_exam_user_unit ON user_exams (user_id, unit_id);