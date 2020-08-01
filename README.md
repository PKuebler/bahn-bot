# Bahn Bot 

German telegram bot as notification in case of train delays. It's use [marudor.de](https://marudor.de) as hafas api.


![GitHub](https://img.shields.io/github/license/pkuebler/bahn-bot?style=for-the-badge)
[![Travis (.org)](https://img.shields.io/travis/pkuebler/bahn-bot?style=for-the-badge)](https://travis-ci.org/github/PKuebler/bahn-bot)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/pkuebler/bahn-bot?style=for-the-badge)](https://hub.docker.com/repository/docker/pkuebler/bahn-bot)
[![MicroBadger Layers](https://img.shields.io/microbadger/layers/pkuebler/bahn-bot?style=for-the-badge)](https://hub.docker.com/repository/docker/pkuebler/bahn-bot)
[![Docker Image Size (tag)](https://img.shields.io/docker/image-size/pkuebler/bahn-bot/latest?style=for-the-badge)](https://hub.docker.com/repository/docker/pkuebler/bahn-bot)
[![Docker Automated build](https://img.shields.io/docker/cloud/automated/pkuebler/bahn-bot?style=for-the-badge)](https://hub.docker.com/repository/docker/pkuebler/bahn-bot)
[![Docker Build Status](https://img.shields.io/docker/cloud/build/pkuebler/bahn-bot?style=for-the-badge)](https://hub.docker.com/repository/docker/pkuebler/bahn-bot)

## Features

- Beobachte Züge
- Individuelle Schwellwerte

## Commands

- `/help` Befehlsübersicht
- `/myalarms` Alle gesetzten Verspätungsalarme
- `/newalarm` Neuer Verspätungsalarm
- `/cancel` Abbrechen des aktuellen Vorgangs

## Config

All parameters can be overwritten by ENV variable.

```json
// config.json
{
    "api": {
        "endpoint": "https://marudor.de/api"
    },
    "telegram": {
        "key": ""
    },
    "database": {
        "dialect": "sqlite3",
        "path": ":memory:"
    },
    "loglevel": "trace"
}
```

- `API_ENDPOINT` marudor endpoint
- `TELEGRAM_KEY` telegram bot key from [Telegram BotFahter](https://core.telegram.org/bots#6-botfather)
- `DB_DIALECT` currently only `sqlite3`
- `DB_PATH` database config path.
    - sqlite3
        - `:memory:` inmemory database
        - `./path/to/databasefile` filebased
- `LOG_LEVEL` loglevel `info`, `trace`, `warn` or `error` 

## Docker Setup

```yml
version: '3'
services:
  bot:
    image: "pkuebler/bahn-bot:latest"
    environment:
        - API_ENDPOINT=https://marudor.de/api
        - TELEGRAM_KEY=
        - DB_DIALECT=sqlite3
        - "DB_PATH=:memory:"
        - LOG_LEVEL=info
```

## DSGVO

- Speichert:
    - ChatID + Zug für den Verspätungsalarm. Wird 2 Tage nach Ankunft des Zuges gelöscht (um Verspätungen abzufangen).
    - ChatID + Aktuelle Operation mit Bot (z.B. newalarm, savealarm, ...)

## ToDo

- Add interface > telegram tests
- Delete old chat states -> notify after delete
- SQL Database Repository tests

## License

MIT