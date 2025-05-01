CREATE TABLE IF NOT EXISTS difficulties (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(10) DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'archived')),
    description_short VARCHAR(200),
    description_long VARCHAR(3000),
    total_categories INTEGER[] NOT NULL DEFAULT '{}',
    image_url VARCHAR(255),
    update_news JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS difficulties_news (
    id SERIAL PRIMARY KEY,
    difficulty_id INTEGER NOT NULL REFERENCES difficulties(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
