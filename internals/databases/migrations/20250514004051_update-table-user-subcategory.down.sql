-- Hapus index current_version jika ada
DROP INDEX IF EXISTS idx_user_subcategory_version;

-- Hapus kolom current_version dari tabel user_subcategory
ALTER TABLE user_subcategory DROP COLUMN IF EXISTS current_version;
