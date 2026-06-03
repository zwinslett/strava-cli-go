// Package model defines the data structure of the Strava API responses
package model

type Gear struct {
	ID                string  `json:"id"`
	Primary           bool    `json:"primary"`
	Name              string  `json:"name"`
	Nickname          *string `json:"nickname"`
	ResourceState     int     `json:"resource_state"`
	Retired           bool    `json:"retired"`
	Distance          float64 `json:"distance"`
	ConvertedDistance float64 `json:"converted_distance"`
}
