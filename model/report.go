package model

type ActivityReport struct {
	Activity DetailedActivity                `json:"activity"`
	Zones    []DistributionBucketsAggregated `json:"zones"`
}
