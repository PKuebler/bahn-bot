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
