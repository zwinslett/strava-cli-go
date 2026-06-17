// Package calculator  contains functions for transforming, aggregating, and processing Strava activitiy data.
package calculator

import "github.com/zwinslett/strava-cli-go/model"

func MetersToMiles(meters float64) float64 {
	return meters * 0.000621371
}

func FilterByType(sportType string, activities []model.Activity) []model.Activity {
	var filtered []model.Activity

	for _, activity := range activities {
		if activity.SportType == sportType {
			filtered = append(filtered, activity)
		}
	}
	return filtered
}

func ActivityComparison(previousDistance float64, currentDistance float64, previousMovingTime int, currentMovingTime int, previousAverageHeartrate float64, currentAverageHeartRate float64) (float64, float64, float64) {
	distanceDiff := (currentDistance - previousDistance) / previousDistance * 100
	movingTimeDiff := float64(currentMovingTime-previousMovingTime) / float64(previousMovingTime) * 100
	averageHeartrateDiff := (currentAverageHeartRate - previousAverageHeartrate) / previousAverageHeartrate * 100

	return distanceDiff, movingTimeDiff, averageHeartrateDiff
}
