
# ssd-bot-go
u/_SSD_BOT_

This is the golang port of the original ssd bot for commenting ssd information in the /r/buildapcsales subreddit.

You can find the original bot [here](https://github.com/ocmarin/ssd-bot).

# Running locally
```shell
cp .env.example .env
export $(grep -v '^#' .env | xargs)
go build main.go
./main
```

# Running on docker
```shell
cp .env.example .env
docker build --tag ssd-bot-go .
docker compose up -d
```
# Acknowledgement
Thanks [TechPowerup](https://www.techpowerup.com/ssd-specs/) for providing me their api access to their SSD database!
