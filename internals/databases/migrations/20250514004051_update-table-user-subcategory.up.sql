-- Tambahkan kolom version_current ke tabel user_subcategory
ALTER TABLE user_subcategory
ADD COLUMN current_version INTEGER DEFAULT 1;

-- Buat index agar pencarian berdasarkan current_version lebih efisien
CREATE INDEX IF NOT EXISTS idx_user_subcategory_version ON user_subcategory(current_version);
