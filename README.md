# Etnograbber

Etnograbber proxies the tracks endpoint of the SoundCloud API while
managing the underlying OAuth token lifecycle in order to not exhaust
the weekly/monthly new token emissions from SC (renewals are infinite).

## Why?

Managing tokens is a pain in sub-par languages that are short lived
like Personal Home Page, so it's better to wrap the legacy Personal
Home Page application.

## Develop

* `cp config.dist.toml config.toml`
* `just or go run/build/test`

## Deploy

* `just buildserver`
* copy `bin/linux/etnograbber` to the server
* copy a valid `config.toml`
* run

## Bovino seal of approval:

![](https://upload.wikimedia.org/wikipedia/en/2/21/Blink-182_-_Dude_Ranch_cover.jpg)
