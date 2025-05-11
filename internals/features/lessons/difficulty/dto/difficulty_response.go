package model

type DifficultyResponse struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Status           string `json:"status"`
	DescriptionShort string `json:"description_short"`
	DescriptionLong  string `json:"description_long"`
	TotalCategories  []int  `json:"total_categories"`
	ImageURL         string `json:"image_url"`
}
