package calculator

import "fmt"

func PrettifiedTime(movingTime int) string {
	minutes := movingTime / 60
	seconds := movingTime % 60
	return fmt.Sprintf("%d Minutes %d Seconds", minutes, seconds)
}

func SecsToMins(seconds float64) float64 {
	return seconds / 60
}
