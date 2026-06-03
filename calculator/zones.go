package calculator

import (
	"sort"

	"github.com/zwinslett/strava-cli-go/model"
)

type bucketKey struct {
	min float64
	max float64
}

func AggregateZones(zones []model.Zones, zoneType ZoneType) []model.DistributionBucketsAggregated {
	var filtered []model.Zones
	for _, zone := range zones {
		if zone.Type == string(zoneType) {
			filtered = append(filtered, zone)
		}
	}
	totals := make(map[bucketKey]float64)
	for _, zone := range filtered {
		for _, bucket := range zone.Buckets {
			key := bucketKey{bucket.Min, bucket.Max}
			totals[key] += bucket.Time
		}
	}
	var aggregatedBuckets []model.DistributionBucketsAggregated

	for zone, time := range totals {
		aggregatedBuckets = append(aggregatedBuckets, model.DistributionBucketsAggregated{
			Min:  zone.min,
			Max:  zone.max,
			Time: time,
		})
	}
	var totalTime float64
	for _, bucket := range aggregatedBuckets {
		totalTime += bucket.Time
	}

	for i := range aggregatedBuckets {
		aggregatedBuckets[i].Percent = aggregatedBuckets[i].Time / totalTime
	}
	sort.Slice(aggregatedBuckets, func(i, j int) bool {
		return aggregatedBuckets[i].Min < aggregatedBuckets[j].Min
	})

	return aggregatedBuckets
}
