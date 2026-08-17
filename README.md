# Discord Welcome Bot

A simple Discord bot that welcomes users when they join a server and automatically assigns them a role

## Requirements
- Go 1.20+
- A Discord bot token
- Server Members Intent enabled

## Installation
Clone the repository
```
git clone https://github.com/Waffle-SnD/waffle-bot.git
cd welcome-bot
```

Install the dependencies
```
go mod tidy
```

Run or build the bot
```
go run .
```

```
go build -ldflags "-s -w" .
```
