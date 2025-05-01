package main

import (
	"log"

	"quizku/internals/configs"
	"quizku/internals/seeds"
)

func main() {
	configs.LoadEnv() // <-- penting
	db := configs.InitDB()
	log.Println("🚀 Menjalankan semua seed...")
	seeds.RunAllSeeds(db)
}
