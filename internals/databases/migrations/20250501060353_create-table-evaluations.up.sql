CREATE TABLE IF NOT EXISTS evaluations (
    evaluation_id SERIAL PRIMARY KEY,
    evaluation_name VARCHAR(50) NOT NULL,
    evaluation_status VARCHAR(10) CHECK (evaluation_status IN ('active', 'pending', 'archived')) DEFAULT 'pending',
    evaluation_total_question INTEGER[] NOT NULL DEFAULT '{}',
    evaluation_icon_url VARCHAR(100),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    evaluation_unit_id INT REFERENCES units(unit_id) ON DELETE CASCADE,
    evaluation_created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- Indexing untuk performa
CREATE INDEX IF NOT EXISTS idx_evaluations_status ON evaluations(evaluation_status);
CREATE INDEX IF NOT EXISTS idx_evaluations_unit_id ON evaluations(evaluation_unit_id);
CREATE INDEX IF NOT EXISTS idx_evaluations_created_by ON evaluations(evaluation_created_by);


CREATE TABLE IF NOT EXISTS user_evaluations (
    user_evaluation_id SERIAL PRIMARY KEY,
    user_evaluation_user_id UUID NOT NULL 
        REFERENCES users(id) ON DELETE CASCADE,
    user_evaluation_evaluation_id INTEGER NOT NULL 
        REFERENCES evaluations(evaluation_id) ON DELETE CASCADE,
    user_evaluation_unit_id INTEGER NOT NULL 
        REFERENCES units(unit_id) ON DELETE CASCADE,
    user_evaluation_attempt INTEGER NOT NULL DEFAULT 1,
    user_evaluation_percentage_grade INTEGER NOT NULL DEFAULT 0,
    user_evaluation_time_duration INTEGER NOT NULL DEFAULT 0,
    user_evaluation_point INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS idx_user_eval_eval 
    ON user_evaluations (user_evaluation_user_id, user_evaluation_evaluation_id);

CREATE INDEX IF NOT EXISTS idx_user_eval_unit 
    ON user_evaluations (user_evaluation_user_id, user_evaluation_unit_id);
