-- +migrate Up

CREATE TABLE IF NOT EXISTS user_point_logs (
    user_point_log_id SERIAL PRIMARY KEY,                             -- ID unik log poin
    user_point_log_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- ID user yang mendapat poin
    user_point_log_points INTEGER NOT NULL,                           -- Jumlah poin yang diberikan
    user_point_log_source_type INTEGER NOT NULL,                      -- Tipe sumber poin (0=reading, 1=quiz, 2=evaluation, 3=exam, dst)
    user_point_log_source_id INTEGER,                                 -- ID entitas sumber poin (boleh NULL jika tidak relevan)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP                    -- Timestamp saat poin dicatat
);

-- Index gabungan untuk optimasi pencarian berdasarkan user dan sumber
CREATE INDEX IF NOT EXISTS idx_user_point_logs_user_source
    ON user_point_logs(user_point_log_user_id, user_point_log_source_type, user_point_log_source_id);

-- Index tambahan untuk pencarian umum berdasarkan user_id
CREATE INDEX IF NOT EXISTS idx_user_point_logs_user_id
    ON user_point_logs(user_point_log_user_id);
