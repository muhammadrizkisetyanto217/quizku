-- FINAL: Struktur tabel readings yang telah dioptimasi

CREATE TABLE IF NOT EXISTS readings (
	id SERIAL PRIMARY KEY,
	title VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(10) DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'archived')),
	description_long TEXT NOT NULL,
	tooltips_id INTEGER[] NOT NULL DEFAULT '{}',
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
