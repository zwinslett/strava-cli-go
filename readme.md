# Strava CLI 

This is a Golang (Cobra) command line interface and daemon that interacts with the Strava and Telegram APIs that allows users to: 

- Request data from Strava in a JSON format with composable commands for piping with `jq` etc. 
- Request data from Strava right from their terminal for viewing. 
- Configure a daemon to send custom notifications based on Strava data over the Telegram Bot API. 

It is useful for those who want a minimal and customizable interface for reading and analyzing Strava data, those who want to automate notifications and other processes based on their activities, and those who want a lightweight (i.e. not a full MCP) interface that an LLM can interact with. 

> Note that this CLI is designed to only work with running activities. It should be relatively easy to modify it to work with other activities, but that is not currently on my roadmap. Feel free to make those additions yourself. You'll find that activities are filtered by the `FilterByType` function inside the calculator package.  

## Commands 
|  Command | Flags  |  Use / Subcommands |
|---|---|---|
|  `last` | `--json`  | Retrieve data about the last activity.  |
| `stats`  |  `--json`, `--weekly`, `--monthly` | Retrieve *rolling* weekly or monthly stats.  |
| `zones`  |  `--json`, `--weekly`, `--monthly` | Retrieve zones data from a *rolling* weekly or monthly range.  |
| `notify` | | In conjunction with `stats` and `last` subcommands sends a notification to Telegram about stats or an individual activity. `stats` supports the `--weekly` and `--monthly` flags. |
|`activity`|`--json` | Retrieves data about a specific activity. Requires an activity id as an argument (`int64`)|
| `daemon`| | Kicks off the daemon for scheduling Telegram messages |

## Setting up environment variables

## Setting up a Strava API application 

## Setting up a Telegram Bot 
