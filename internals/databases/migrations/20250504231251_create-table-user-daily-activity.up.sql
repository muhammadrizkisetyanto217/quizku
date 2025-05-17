CREATE TABLE IF NOT EXISTS user_daily_activities (
    user_daily_activity_id SERIAL PRIMARY KEY,
    user_daily_activity_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_daily_activity_activity_date DATE NOT NULL,
    user_daily_activity_amount_day INT NOT NULL DEFAULT 1,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(user_daily_activity_user_id, user_daily_activity_activity_date)
);
