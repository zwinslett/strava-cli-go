# Strava CLI 

This is a Golang (Cobra) command line interface for interacting with the Strava and Telegram APIs that allows users to: 

- Request data from Strava in a JSON format for piping with `jq` etc. 
- Request data from Strava right from their terminal for viewing. 
- Send notifications over the Telegram Bot API to chats. 

It is useful for those who want a minimal and customizable interface for reading and analyzing Strava data, those who want to automate notifications and other processes based on their activities, and those who want a lightweight interface for an LLM. 

## Commands 
|  Command | Flags  |  Use |
|---|---|---|
|  `last` | `--json`  | Retrieve data about the last activity.  |
| `stats`  |  `--json`, `--weekly`, `--monthly` | Retrieve weekly or monthly stats.  |
| `zones`  |  `--json`, `--weekly`, `--monthly` | Retrieve zones data from a weekly or monthly range.  |
| `notify` | | In conjunction with `stats` and `last` subcommands sends a notification to Telegram about stats or an individual activity. `stats` supports the `--weekly` and `--monthly` flags. |
|`activity`| | Retrieves data about a specific activity. Requires an activity id as an argument (`int64`)|
