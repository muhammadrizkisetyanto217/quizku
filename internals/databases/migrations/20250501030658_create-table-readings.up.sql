-- FINAL: Struktur tabel readings yang telah dioptimasi

CREATE TABLE IF NOT EXISTS readings (
	id SERIAL PRIMARY KEY,
	title VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'archived')),
	description_long TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP,
	unit_id INT REFERENCES units(id) ON DELETE CASCADE,
	created_by UUID REFERENCES users(id) ON DELETE CASCADE
);

-- Indexing untuk performa
CREATE INDEX IF NOT EXISTS idx_readings_status ON readings(status);
CREATE INDEX IF NOT EXISTS idx_readings_unit_id ON readings(unit_id);
CREATE INDEX IF NOT EXISTS idx_readings_created_by ON readings(created_by);


CREATE TABLE IF NOT EXISTS user_readings (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    reading_id INTEGER NOT NULL REFERENCES readings(id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    attempt INTEGER DEFAULT 1 NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index
CREATE INDEX IF NOT EXISTS idx_user_readings_user_id_reading_id ON user_readings (user_id, reading_id);
CREATE INDEX IF NOT EXISTS idx_user_readings_user_id_unit_id ON user_readings (user_id, unit_id);