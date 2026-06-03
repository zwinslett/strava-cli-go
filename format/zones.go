package format

import (
	"fmt"

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
