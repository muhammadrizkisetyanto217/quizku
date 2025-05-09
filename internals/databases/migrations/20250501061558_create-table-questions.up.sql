CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    question_text TEXT NOT NULL,
    question_answer TEXT[] NOT NULL,  -- array jawaban pilihan (["A", "B", "C", "D"])
    question_correct TEXT NOT NULL CHECK (char_length(question_correct) <= 50),
    paragraph_help TEXT NOT NULL,
    explain_question TEXT NOT NULL,
    answer_text TEXT NOT NULL,
    tooltips_id INT[] DEFAULT '{}',  -- opsional relasi ke tabel tooltips
    donation_id INT REFERENCES user_question_donations(id) ON DELETE SET NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'archived')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
-- Index status untuk filter cepat soal aktif/pending
CREATE INDEX IF NOT EXISTS idx_questions_status ON questions(status);
-- Index tooltips_id untuk query berdasarkan isi array (jika digunakan)
CREATE INDEX IF NOT EXISTS idx_questions_tooltips_id ON questions USING GIN(tooltips_id);


CREATE TABLE IF NOT EXISTS question_links (
    id SERIAL PRIMARY KEY,
    question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    target_type SMALLINT NOT NULL CHECK (target_type IN (1, 2, 3, 4)),
    target_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_question_links_question_id ON question_links(question_id);
CREATE INDEX IF NOT EXISTS idx_question_links_target ON question_links(target_type, target_id);


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
CREATE INDEX IF NOT EXISTS idx_user_questions_user_id ON user_questions (user_id);
CREATE INDEX IF NOT EXISTS idx_user_questions_question_id ON user_questions (question_id);
CREATE INDEX IF NOT EXISTS idx_user_questions_source_type_id ON user_questions (source_type_id);
CREATE INDEX IF NOT EXISTS idx_user_questions_source_id ON user_questions (source_id);
