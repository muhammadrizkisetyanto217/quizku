-- +migrate Up
CREATE TABLE IF NOT EXISTS user_progress (
    user_progress_id SERIAL PRIMARY KEY,
    user_progress_user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    user_progress_total_points INTEGER NOT NULL DEFAULT 0,
    user_progress_level INTEGER NOT NULL DEFAULT 1,
    user_progress_rank INTEGER NOT NULL DEFAULT 1,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk lookup cepat
CREATE INDEX IF NOT EXISTS idx_user_progress_user_id 
    ON user_progress(user_progress_user_id);
