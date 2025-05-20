package model

type DifficultyResponse struct {
	DifficultyID               uint   `json:"difficulty_id"`
	DifficultyName             string `json:"difficulty_name"`
	DifficultyStatus           string `json:"difficulty_status"`
	DifficultyDescriptionShort string `json:"difficulty_description_short"`
	DifficultyDescriptionLong  string `json:"difficulty_description_long"`
	DifficultyTotalCategories  []int  `json:"difficulty_total_categories"`
	DifficultyImageURL         string `json:"difficulty_image_url"`
}
