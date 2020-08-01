# Bahn Bot ![Travis CI](https://api.travis-ci.org/PKuebler/bahn-bot.svg?branch=master&status=started)

German telegram bot as notification in case of train delays. It's use [marudor.de](https://marudor.de) as hafas api.

## Features

- Beobachte Züge
- Individuelle Schwellwerte

## Commands

- `/help` Befehlsübersicht
- `/myalarms` Alle gesetzten Verspätungsalarme
- `/newalarm` Neuer Verspätungsalarm
- `/cancel` Abbrechen des aktuellen Vorgangs

## Config

> todo: config

## Docker Setup

> todo: docker-compose

## DSGVO

- Speichert:
    - ChatID + Zug für den Verspätungsalarm. Wird 2 Tage nach Ankunft des Zuges gelöscht (um Verspätungen abzufangen).
    - ChatID + Aktuelle Operation mit Bot (z.B. newalarm, savealarm, ...)

## ToDo

- Add interface > telegram tests
- Delete old chat states -> notify after delete
- SQL Database Repository tests
