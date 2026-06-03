// Package model defines the data structure of the Strava API responses
package model

type Activity struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Distance         float64 `json:"distance"`
	MovingTime       int     `json:"moving_time"`
	AverageHeartrate float64 `json:"average_heartrate"`
	MaxHeartrate     float64 `json:"max_heartrate"`
	SufferScore      float64 `json:"suffer_score"`
	Calories         float64 `json:"calories"`
	Description      string  `json:"description"`
}

type DetailedActivity struct {
	Activity
	SplitsStandard []Split `json:"splits_standard"`
	Gear           Gear    `json:"gear"`
}
