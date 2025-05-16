-- +migrate Up
CREATE TABLE IF NOT EXISTS tooltips (
    tooltip_id SERIAL PRIMARY KEY,                                         -- ID unik untuk setiap tooltip
    tooltip_keyword TEXT NOT NULL UNIQUE CHECK (char_length(tooltip_keyword) <= 100),  
        -- Kata kunci atau topik tooltip, harus unik dan maksimal 100 karakter
    tooltip_description_short TEXT NOT NULL CHECK (char_length(tooltip_description_short) <= 200),  
        -- Ringkasan singkat tooltip (untuk pratinjau), maksimal 200 karakter
    tooltip_description_long TEXT NOT NULL,                                -- Penjelasan panjang/detail untuk edukasi pengguna
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,              -- Waktu pertama kali data dibuat
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP               -- Waktu terakhir data diupdate
);

-- Index tambahan untuk optimasi pencarian tooltip berdasarkan keyword, terutama untuk query LIKE
CREATE INDEX IF NOT EXISTS idx_tooltips_keyword_lower 
    ON tooltips (LOWER(tooltip_keyword));
