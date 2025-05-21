CREATE TABLE IF NOT EXISTS questions (
    question_id SERIAL PRIMARY KEY,  -- ID unik untuk setiap pertanyaan
    question_text TEXT NOT NULL,  -- Isi teks utama dari soal/pertanyaan
    question_answer_choices TEXT[] NOT NULL,  
        -- Opsi jawaban pilihan ganda (bisa lebih dari 2)
    question_correct_answer TEXT NOT NULL CHECK (char_length(question_correct_answer) <= 50),  
        -- Jawaban yang benar dari pilihan (harus sesuai salah satu dari array question_answer_choices)
    question_paragraph_help TEXT NOT NULL,  
        -- Teks bacaan pendukung atau paragraf referensi
    question_explanation TEXT NOT NULL,  
        -- Penjelasan mengapa jawaban tersebut benar
    question_answer_text TEXT NOT NULL,  
        -- Jawaban dalam bentuk teks panjang (jika diperlukan)
    question_status VARCHAR(10) NOT NULL DEFAULT 'pending' 
        CHECK (question_status IN ('active', 'pending', 'archived')),  
        -- Status pertanyaan: aktif, pending, atau diarsipkan
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Waktu dibuat
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- Waktu diperbarui
    deleted_at TIMESTAMP  -- Soft delete (jika ada)
);

-- Index untuk mempercepat pencarian berdasarkan status
CREATE INDEX IF NOT EXISTS idx_questions_status ON questions(question_status);


CREATE TABLE IF NOT EXISTS question_links (
    question_link_id SERIAL PRIMARY KEY,
    question_link_question_id INT NOT NULL REFERENCES questions(question_id) ON DELETE CASCADE,
    question_link_target_type SMALLINT NOT NULL CHECK (question_link_target_type IN (1, 2, 3, 4)),
    question_link_target_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat pencarian berdasarkan pertanyaan
CREATE INDEX IF NOT EXISTS idx_question_links_question_id 
ON question_links(question_link_question_id);

-- Index untuk mempercepat pencarian berdasarkan target_type dan target_id
CREATE INDEX IF NOT EXISTS idx_question_links_target 
ON question_links(question_link_target_type, question_link_target_id);



CREATE TABLE IF NOT EXISTS question_saved (
    question_saved_id SERIAL PRIMARY KEY,
    question_saved_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_saved_source_type_id INTEGER NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    question_saved_question_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat query pencarian
CREATE INDEX IF NOT EXISTS idx_question_saved_user 
ON question_saved(question_saved_user_id);

CREATE INDEX IF NOT EXISTS idx_question_saved_question 
ON question_saved(question_saved_question_id);


CREATE TABLE IF NOT EXISTS question_mistakes (
    question_mistake_id SERIAL PRIMARY KEY,
    question_mistake_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_mistake_source_type_id INTEGER NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    question_mistake_question_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa pencarian
CREATE INDEX IF NOT EXISTS idx_question_mistakes_user_id 
    ON question_mistakes(question_mistake_user_id);

CREATE INDEX IF NOT EXISTS idx_question_mistakes_question_id 
    ON question_mistakes(question_mistake_question_id);




CREATE TABLE IF NOT EXISTS user_questions (
    user_question_id SERIAL PRIMARY KEY,
    user_question_user_id UUID NOT NULL,
    user_question_question_id INT NOT NULL REFERENCES questions(question_id) ON DELETE CASCADE,
    user_question_selected_answer TEXT NOT NULL,
    user_question_is_correct BOOLEAN NOT NULL,
    user_question_source_type_id INT NOT NULL, -- 1 = Quiz, 2 = Evaluation, 3 = Exam
    user_question_source_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing for performance
CREATE INDEX IF NOT EXISTS idx_user_questions_user_id 
    ON user_questions (user_question_user_id);

CREATE INDEX IF NOT EXISTS idx_user_questions_question_id 
    ON user_questions (user_question_question_id);

CREATE INDEX IF NOT EXISTS idx_user_questions_source_type_id 
    ON user_questions (user_question_source_type_id);

CREATE INDEX IF NOT EXISTS idx_user_questions_source_id 
    ON user_questions (user_question_source_id);
