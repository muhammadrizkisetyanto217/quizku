-- FINAL: Struktur tabel readings yang telah dioptimasi

-- FINAL: Struktur tabel readings yang telah di-refactor semantik

CREATE TABLE IF NOT EXISTS readings (
    reading_id SERIAL PRIMARY KEY,  -- ID unik bacaan
    reading_title VARCHAR(50) UNIQUE NOT NULL,  -- Judul bacaan
    reading_status VARCHAR(10) DEFAULT 'pending' 
        CHECK (reading_status IN ('active', 'pending', 'archived')),  -- Status: active/pending/archived
    reading_description_long TEXT NOT NULL,  -- Deskripsi panjang bacaan
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp dibuat
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Timestamp diperbarui
    deleted_at TIMESTAMP,  -- Soft delete
    reading_unit_id INT REFERENCES units(id) ON DELETE CASCADE,  -- Relasi ke unit
    reading_created_by UUID REFERENCES users(id) ON DELETE CASCADE  -- Relasi ke user pembuat
);

-- 🔍 Indexing untuk performa query
CREATE INDEX IF NOT EXISTS idx_readings_status ON readings(reading_status);
CREATE INDEX IF NOT EXISTS idx_readings_unit_id ON readings(reading_unit_id);
CREATE INDEX IF NOT EXISTS idx_readings_created_by ON readings(reading_created_by);



-- +migrate Up
CREATE TABLE IF NOT EXISTS user_readings (
    user_reading_id SERIAL PRIMARY KEY,
    user_reading_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_reading_reading_id INTEGER NOT NULL REFERENCES readings(reading_id) ON DELETE CASCADE,
    user_reading_unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    user_reading_attempt INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ✅ Indexing untuk performa dan pencarian gabungan
CREATE INDEX IF NOT EXISTS idx_user_readings_user_id_reading_id 
    ON user_readings (user_reading_user_id, user_reading_reading_id);

CREATE INDEX IF NOT EXISTS idx_user_readings_user_id_unit_id 
    ON user_readings (user_reading_user_id, user_reading_unit_id);
