package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
func ConnectDB() {
	dsn := "postgresql://postgres.iqtqixqqsubcjbqdmgsw:Wedangjahe217@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres"

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // ✅ Ini kunci utama
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Gagal koneksi ke Supabase:", err)
	}

	DB = db
	fmt.Println("🚀 Berhasil konek ke Supabase PostgreSQL sekarang!")
}
