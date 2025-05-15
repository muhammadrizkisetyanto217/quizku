-- SQL untuk tabel issued_certificates
CREATE TABLE issued_certificates (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    subcategory_id INTEGER NOT NULL,
    certificate_version_id INTEGER NOT NULL,
    version_current INTEGER NOT NULL,
    is_up_to_date BOOLEAN NOT NULL DEFAULT true,
    slug_url TEXT UNIQUE NOT NULL,
    issued_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE certificate_versions (
    id SERIAL PRIMARY KEY,
    subcategory_id INTEGER NOT NULL REFERENCES subcategories(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    total_themes INTEGER NOT NULL DEFAULT 0,
    note TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE
);