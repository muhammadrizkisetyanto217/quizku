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
