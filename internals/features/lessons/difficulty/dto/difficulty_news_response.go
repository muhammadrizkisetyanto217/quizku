package dto

import "time"

type DifficultyNewsDTO struct {
	DifficultyNewsID          uint      `json:"difficulty_news_id"`
	DifficultyNewsTitle       string    `json:"difficulty_news_title"`
	DifficultyNewsDescription string    `json:"difficulty_news_description"`
	CreatedAt                 time.Time `json:"created_at"`
}
