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
