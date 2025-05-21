-- ✅ TABLE: subcategories
CREATE TABLE IF NOT EXISTS subcategories (
    subcategory_id SERIAL PRIMARY KEY,
    subcategory_name VARCHAR(255) NOT NULL,
    subcategory_status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (subcategory_status IN ('active', 'pending', 'archived')),
    subcategory_description_long TEXT,
    subcategory_total_themes_or_levels INTEGER[] NOT NULL DEFAULT '{}',
    subcategory_image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    subcategory_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    subcategory_deleted_at TIMESTAMP,
    subcategory_category_id INT NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,

    CONSTRAINT unique_subcategory_name_per_category UNIQUE (subcategory_name, subcategory_category_id)
);

-- ✅ Index untuk performa query
CREATE INDEX IF NOT EXISTS idx_subcategory_status 
    ON subcategories(subcategory_status);

CREATE INDEX IF NOT EXISTS idx_subcategory_category 
    ON subcategories(subcategory_category_id);

CREATE INDEX IF NOT EXISTS idx_subcategory_cat_status 
    ON subcategories(subcategory_category_id, subcategory_status);


-- ✅ TABLE: subcategories_news
-- ✅ TABLE: subcategories_news (refactored semantik)
CREATE TABLE IF NOT EXISTS subcategories_news (
    subcategory_news_id SERIAL PRIMARY KEY,
    subcategory_news_subcategory_id INTEGER NOT NULL REFERENCES subcategories(subcategory_id) ON DELETE CASCADE,
    subcategory_news_title VARCHAR(255) NOT NULL,
    subcategory_news_description TEXT NOT NULL,
    subcategory_news_is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- ✅ Index untuk performa pencarian
CREATE INDEX IF NOT EXISTS idx_subcategory_news_subcategory_id 
    ON subcategories_news (subcategory_news_subcategory_id);

CREATE INDEX IF NOT EXISTS idx_subcategory_news_is_public 
    ON subcategories_news (subcategory_news_is_public);

CREATE INDEX IF NOT EXISTS idx_subcategory_news_subcat_public 
    ON subcategories_news (subcategory_news_subcategory_id, subcategory_news_is_public);



-- ✅ TABLE: user_subcategory (refactored + created_at/updated_at tetap)
CREATE TABLE IF NOT EXISTS user_subcategories (
    user_subcategory_id SERIAL PRIMARY KEY,
    user_subcategory_user_id UUID NOT NULL,user_subcategory_subcategory_id INTEGER NOT NULL REFERENCES subcategories(subcategory_id) ON DELETE CASCADE,
    user_subcategory_complete_themes_or_levels JSONB NOT NULL DEFAULT '{}'::jsonb,
    user_subcategory_grade_result INTEGER NOT NULL DEFAULT 0,
    user_subcategory_current_version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Index untuk performa
CREATE INDEX IF NOT EXISTS idx_user_subcategory_user_id 
    ON user_subcategories (user_subcategory_user_id);

CREATE INDEX IF NOT EXISTS idx_user_subcategory_subcategory_id 
    ON user_subcategories (user_subcategory_subcategory_id);

CREATE INDEX IF NOT EXISTS idx_user_subcategory_user_subcategory 
    ON user_subcategories(user_subcategory_user_id, user_subcategory_subcategory_id);
