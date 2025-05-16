CREATE TABLE IF NOT EXISTS level_requirements (
    level_req_id SERIAL PRIMARY KEY,                                -- ID unik untuk level requirement
    level_req_level INT UNIQUE NOT NULL,                            -- Level numerik (misalnya: 1, 2, 3)
    level_req_name VARCHAR(100),                                    -- Nama atau label level (misal: Pemula, Mahir)
    level_req_min_points INT NOT NULL,                              -- Minimum poin yang dibutuhkan untuk mencapai level ini
    level_req_max_points INT,                                       -- Maksimum poin (jika ada), NULL artinya tak terbatas
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                 -- Waktu pembuatan data
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP                  -- Waktu terakhir update data
);

CREATE TABLE IF NOT EXISTS rank_requirements (
    rank_req_id SERIAL PRIMARY KEY,                                 -- ID unik untuk rank requirement
    rank_req_rank INTEGER UNIQUE NOT NULL,                          -- Nomor urut pangkat (misal: 1, 2, 3)
    rank_req_name VARCHAR(100),                                     -- Nama atau label pangkat (misal: Prajurit, Jenderal)
    rank_req_min_level INTEGER NOT NULL,                            -- Level minimum yang diperlukan untuk mencapai pangkat ini
    rank_req_max_level INTEGER,                                     -- Level maksimum (boleh NULL jika tidak terbatas)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,        -- Timestamp pembuatan data
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP         -- Timestamp pembaruan data
);
