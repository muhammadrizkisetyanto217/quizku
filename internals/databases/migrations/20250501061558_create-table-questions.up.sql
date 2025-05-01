CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    source_type_id INT NOT NULL CHECK (source_type_id IN (1, 2, 3)), -- enum enforcement
    source_id INT NOT NULL,
    question_text TEXT NOT NULL, -- ubah dari VARCHAR(200) ke TEXT
    question_answer TEXT[] NOT NULL,
    question_correct TEXT NOT NULL CHECK (char_length(question_correct) <= 50),
    tooltips_id INT[] DEFAULT '{}',
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'archived')),
    paragraph_help TEXT NOT NULL,
    explain_question TEXT NOT NULL,
    answer_text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Index untuk performa query
CREATE INDEX IF NOT EXISTS idx_questions_source ON questions(source_type_id, source_id);
CREATE INDEX IF NOT EXISTS idx_questions_status ON questions(status);
