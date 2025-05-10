CREATE TABLE donations (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    amount INTEGER NOT NULL,
    message TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, paid, expired, canceled
    order_id VARCHAR(100) UNIQUE NOT NULL,
    payment_token TEXT,
    payment_gateway VARCHAR(50) DEFAULT 'midtrans',
    payment_method VARCHAR(50),
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);


CREATE TABLE IF NOT EXISTS donation_questions (
    id SERIAL PRIMARY KEY,
    donation_id BIGINT REFERENCES donations(id) ON DELETE CASCADE,
    question_id BIGINT REFERENCES questions(id) ON DELETE CASCADE,
    user_progress_id BIGINT REFERENCES user_progress(id) ON DELETE SET NULL,
    user_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
