package model

type Zones struct {
	Buckets []DistributionBuckets `json:"distribution_buckets"`
	Type    string                `json:"type"`
}

type DistributionBuckets struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max"`
	Time float64 `json:"time"`
}

type DistributionBucketsAggregated struct {
	Min     float64
	Max     float64
	Time    float64
	Percent float64
}
