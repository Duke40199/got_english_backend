package utils

import (
	"fmt"

	"github.com/golang/got_english_backend/models"
)

func CalculateAverageRating(ratingList *[]models.Rating) float32 {
	var ratingSum float32
	for i := 0; i < len(*ratingList); i++ {
		ratingSum += ((*ratingList)[i]).Score
	}
	return ratingSum / float32(len(*ratingList))
}

func CalculateExpertWeightedRating(expertID uint, ratingList *[]models.Rating, ratingAlgorithm *models.RatingAlgorithm) float32 {
	//number of votes given expert is rated
	var numberOfExpertRatings float32
	//minimum votes required, tunable parameter.
	var minimumRatingCount = ratingAlgorithm.MinimumRatingCount
	//average rating of given expert
	var expertAverageRating float32
	//average rating of all expert
	var averageAllExpertRating = ratingAlgorithm.AverageAllExpertsRating
	//
	var weightedRating float32
	//rating of session that expert is in
	var expertSessionRating []models.Rating
	if len(*ratingList) > 0 {
		for i := 0; i < len(*ratingList); i++ {
			if ((*ratingList)[i]).TranslationSession != nil && ((*ratingList)[i]).TranslationSession.Expert.ID == expertID {
				expertSessionRating = append(expertSessionRating, (*ratingList)[i])
			}
			if ((*ratingList)[i]).LiveCallSession != nil && ((*ratingList)[i]).LiveCallSession.Expert.ID == expertID {
				expertSessionRating = append(expertSessionRating, (*ratingList)[i])
			}
			if ((*ratingList)[i]).MessagingSession != nil && ((*ratingList)[i]).MessagingSession.Expert.ID == expertID {
				expertSessionRating = append(expertSessionRating, (*ratingList)[i])
			}
		}
	}
	if uint(len((expertSessionRating))) < minimumRatingCount {
		return 0
	}
	expertAverageRating = CalculateAverageRating(&expertSessionRating)
	numberOfExpertRatings = float32(len(expertSessionRating))
	fmt.Printf("==== expertratingCount:%f\n", numberOfExpertRatings)
	weightedRating =
		((numberOfExpertRatings / (numberOfExpertRatings + float32(minimumRatingCount))) * expertAverageRating) +
			((float32(minimumRatingCount) / (numberOfExpertRatings + float32(minimumRatingCount))) * averageAllExpertRating)
	fmt.Printf("=================weightedrating:%f\n", weightedRating)
	return weightedRating
}
