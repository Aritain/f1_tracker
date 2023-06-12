# F1 Tracker bot

## General info

This is a Telegram bot that is capable of fetching some information regarding F1 season. For now following options are available: driving standing, team standing and next race information

## Installation

To get the bot running simply clone this repo and follow these steps:

1. Create .env file with "TG_TOKEN" variables inside
2. Adjust `TZ_OFFSET` variable under `docker-compose.yml` (by default time is fetched in UTC)
3. Launch the application with `docker-compose up -d`
