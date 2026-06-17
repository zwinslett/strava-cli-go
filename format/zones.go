package format

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/model"
)

func ZonesTable(zones []model.Zones, zoneType calculator.ZoneType) string {
	buckets := calculator.AggregateZones(zones, zoneType)
	table := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("208"))).
		Headers("Zone", "Time", "Percent")
	for _, bucket := range buckets {
		table.Row(
			fmt.Sprintf("%d - %d", int(bucket.Min), int(bucket.Max)),
			calculator.PrettifiedTime(int(bucket.Time)),
			fmt.Sprintf("%.2f", bucket.Percent),
		)
	}
	return table.Render()
}

func ZonesMessage(zones []model.Zones, zoneType calculator.ZoneType) string {
	buckets := calculator.AggregateZones(zones, zoneType)
	var message strings.Builder
	for _, bucket := range buckets {
		fmt.Fprintf(&message,
			"💞 %dbpm - %dbpm\n⏱ %s\n📊 %.2f%%\n\n",
			int(bucket.Min),
			int(bucket.Max),
			calculator.PrettifiedTime(int(bucket.Time)),
			bucket.Percent*100)
	}
	return "📈 <u><b>Zones</b></u>\n" + message.String()
}

func ZonesComparisonMessage(previousZones []model.Zones, currentZones []model.Zones, zoneType calculator.ZoneType) string {
	previousBuckets := calculator.AggregateZones(previousZones, zoneType)
	currentBuckets := calculator.AggregateZones(currentZones, zoneType)
	var message strings.Builder
	for i, bucket := range currentBuckets {
		fmt.Fprintf(&message,
			"💞%dbpm - %dbpm\n⏱ %s\n📊 %.2f%%(%+.2f%%)\n\n",
			int(bucket.Min),
			int(bucket.Max),
			calculator.PrettifiedTime(int(bucket.Time)),
			bucket.Percent*100,
			calculator.CompareBuckets(previousBuckets[i], bucket))
	}
	return "📈 <u><b>Zones</b></u>\n" + message.String()
}
