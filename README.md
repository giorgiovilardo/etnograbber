# Etnograbber

Etnograbber proxies the tracks endpoint of the SoundCloud API while
managing the underlying OAuth token lifecycle in order to not exhaust
the weekly/monthly new token emissions from SC (renewals are infinite).

## Develop

* `cp config.dist.toml config.toml`
* `just or go run/build/test`

## Deploy

* `just buildserver`
* copy `bin/linux/etnograbber` to the server
* copy a valid `config.toml`
* run

Bovino seal of approval:

![](https://upload.wikimedia.org/wikipedia/en/2/21/Blink-182_-_Dude_Ranch_cover.jpg)
