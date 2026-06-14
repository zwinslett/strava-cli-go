package calculator

import (
	"sort"

	"github.com/zwinslett/strava-cli-go/model"
)

func AggregateGear(activities []model.DetailedActivity) []model.Gear {
	totals := make(map[string]float64)
	var aggregatedGear []model.Gear

	for _, activity := range activities {
		key := activity.Gear.Name
		totals[key] += MetersToMiles(activity.Distance)
	}

	for name, distance := range totals {
		aggregatedGear = append(aggregatedGear, model.Gear{
			Name:     name,
			Distance: distance,
		})
	}

	sort.Slice(aggregatedGear, func(i, j int) bool {
		return aggregatedGear[i].Distance > aggregatedGear[j].Distance
	})

	return aggregatedGear
}
