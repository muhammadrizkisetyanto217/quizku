CREATE TABLE IF NOT EXISTS evaluations (
	id SERIAL PRIMARY KEY,
	name_evaluation VARCHAR(50) NOT NULL,
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
CREATE INDEX IF NOT EXISTS idx_evaluations_status ON evaluations(status);
CREATE INDEX IF NOT EXISTS idx_evaluations_unit_id ON evaluations(unit_id);
CREATE INDEX IF NOT EXISTS idx_evaluations_created_by ON evaluations(created_by);


CREATE TABLE IF NOT EXISTS user_evaluations (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    evaluation_id INTEGER NOT NULL REFERENCES evaluations(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    attempt INTEGER DEFAULT 1 NOT NULL,
    percentage_grade INTEGER DEFAULT 0 NOT NULL,
    time_duration INTEGER DEFAULT 0 NOT NULL,
    point INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_eval_user_eval ON user_evaluations (user_id, evaluation_id);
CREATE INDEX IF NOT EXISTS idx_eval_user_unit ON user_evaluations (user_id, unit_id);