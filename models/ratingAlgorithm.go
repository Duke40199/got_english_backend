package models

// Expert model struct
type RatingAlgorithm struct {
	ID                      uint    `gorm:"column:id;not null;unique; primaryKey;" json:"id"`
	MinimumRatingCount      uint    `gorm:"column:minimum_rating_count;default:0" json:"minimum_rating_count"`
	AverageAllExpertsRating float32 `gorm:"column:average_all_experts_rating;default:0" json:"average_all_experts_rating"`
}
