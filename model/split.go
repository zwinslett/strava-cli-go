// Package model defines the data structure of the Strava API responses
package model

type Split struct {
	Distance                  float64 `json:"distance"`
	ElapsedTime               int     `json:"elapsed_time"`
	ElevationDifference       float64 `json:"elevation_difference"`
	MovingTime                int     `json:"moving_time"`
	SplitNumber               int     `json:"split"`
	AverageSpeed              float64 `json:"average_speed"`
	AverageGradeAdjustedSpeed float64 `json:"average_grade_adjusted_speed"`
	AverageHeartrate          float64 `json:"average_heartrate"`
	PaceZone                  int     `json:"pace_zone"`
}
