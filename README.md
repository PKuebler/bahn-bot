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
        "dialect": "postgres",
        "path": "host=myhost port=myport user=gorm dbname=gorm password=mypassword"
    },
    "loglevel": "trace",
    "metrics": false
}
```

- `API_ENDPOINT` marudor endpoint
- `TELEGRAM_KEY` telegram bot key from [Telegram BotFahter](https://core.telegram.org/bots#6-botfather)
- `DB_DIALECT` currently only `mysql` or `postgres`
- `DB_PATH` database config path.
    - mysql
        - `user:password@/dbname?charset=utf8&parseTime=True&loc=Local`
    - postgres
        - `sslmode=disable host=myhost port=myport user=gorm dbname=gorm password=mypassword`
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
        - DB_DIALECT=mysql
        - "DB_PATH=user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
        - LOG_LEVEL=info
```

## Prometheus Endpoint

The Prometheus endpoint is located under `:8080/metrics` when `metrics` is enabled in the Config.

## DSGVO

- Speichert:
    - ChatID + Zug für den Verspätungsalarm. Wird 2 Tage nach Ankunft des Zuges gelöscht (um Verspätungen abzufangen).
    - ChatID + Aktuelle Operation mit Bot (z.B. newalarm, savealarm, ...). Wird 4 Tage nach letzter Interaktion gelöscht.
    - Metrics erfassen verwendung von Zugnummern ohne Verknüpfungen zu Personen.

## ToDo

- Add interface > telegram tests
- SQL Database Repository tests

## License

MIT