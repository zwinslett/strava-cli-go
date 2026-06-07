package format

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/model"
)

func SplitTable(activity model.DetailedActivity) string {
	table := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("208"))).
		Headers("Split", "Duration", "Average Heartrate")

	for _, split := range activity.SplitsStandard {
		table.Row(
			fmt.Sprintf("%d", split.SplitNumber),
			calculator.PrettifiedTime(split.MovingTime),
			fmt.Sprintf("%.2f", split.AverageHeartrate),
		)
	}

	return table.Render()
}

func SplitMessage(activity model.DetailedActivity) string {
	var message strings.Builder
	for _, split := range activity.SplitsStandard {
		fmt.Fprintf(
			&message,
			"🏁 Split: %d\n ⏱ %s\n 💓 %.2fbpm\n\n",
			split.SplitNumber,
			calculator.PrettifiedTime(split.MovingTime),
			split.AverageHeartrate,
		)
	}
	return "📏 <u><b>Splits</b></u>\n" + message.String()
}
