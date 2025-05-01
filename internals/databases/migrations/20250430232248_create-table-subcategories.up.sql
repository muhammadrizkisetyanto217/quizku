CREATE TABLE IF NOT EXISTS subcategories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'archived')),
    description_long TEXT,
    total_themes_or_levels INTEGER[] NOT NULL DEFAULT '{}',
    image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    categories_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE
);

-- Index untuk performa query
CREATE INDEX IF NOT EXISTS idx_subcategories_status ON subcategories(status);
CREATE INDEX IF NOT EXISTS idx_subcategories_category ON subcategories(categories_id);

CREATE TABLE IF NOT EXISTS subcategories_news (
    id SERIAL PRIMARY KEY,
    subcategory_id INTEGER NOT NULL REFERENCES subcategories(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Index untuk performa pencarian per subkategori & status
CREATE INDEX IF NOT EXISTS idx_subnews_subcat ON subcategories_news(subcategory_id);
CREATE INDEX IF NOT EXISTS idx_subnews_is_public ON subcategories_news(is_public);
CREATE INDEX IF NOT EXISTS idx_subnews_public_per_subcat ON subcategories_news(subcategory_id, is_public);
