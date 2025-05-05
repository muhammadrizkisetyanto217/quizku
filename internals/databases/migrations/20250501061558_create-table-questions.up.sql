CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    source_type_id INT NOT NULL CHECK (source_type_id IN (1, 2, 3)), -- enum enforcement
    source_id INT NOT NULL,
    question_text TEXT NOT NULL, -- ubah dari VARCHAR(200) ke TEXT
    question_answer a NOT NULL,
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


CREATE TABLE IF NOT EXISTS question_saved (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    source_type_id INTEGER NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    question_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat query pencarian
CREATE INDEX IF NOT EXISTS idx_question_saved_user ON question_saved(user_id);
CREATE INDEX IF NOT EXISTS idx_question_saved_question ON question_saved(question_id);


CREATE TABLE IF NOT EXISTS question_mistakes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_type_id INTEGER NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    question_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Index untuk performa pencarian
CREATE INDEX IF NOT EXISTS idx_question_mistakes_user_id ON question_mistakes(user_id);
CREATE INDEX IF NOT EXISTS idx_question_mistakes_question_id ON question_mistakes(question_id);



CREATE TABLE IF NOT EXISTS user_questions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    question_id INT NOT NULL,
    selected_answer TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    source_type_id INT NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    source_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing
CREATE INDEX idx_user_questions_user_id ON user_questions (user_id);
CREATE INDEX idx_user_questions_question_id ON user_questions (question_id);
CREATE INDEX idx_user_questions_source_type_id ON user_questions (source_type_id);
CREATE INDEX idx_user_questions_source_id ON user_questions (source_id);
