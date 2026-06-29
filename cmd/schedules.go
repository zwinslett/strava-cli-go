package cmd

type Schedule = string

const (
	DailySchedule   Schedule = "0 14 * * *"
	WeeklySchedule  Schedule = "0 11 * * 1"
	MonthlySchedule Schedule = "0 14 1 * *"
)
