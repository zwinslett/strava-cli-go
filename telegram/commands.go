package telegram

type Command = string

const (
	CmdLatest  Command = "/latest"
	CmdWeekly  Command = "/weekly"
	CmdMonthly Command = "/monthly"
)
