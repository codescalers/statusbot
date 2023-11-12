# Statusbot

## Overview

The Status Bot is a tool designed to send a reminder at specific time

## Getting Started

### Prerequisites

Ensure that you have installed:

- Go programming language (version 1.21 or higher if not using docker)
- Git

### How to start

1. Clone this repository to your local machine:

```bash
git clone https://github.com/codescalers/statusbot
cd statusbot 
```

2. Setup your telegram bot and your env

- Create a new [telegram bot](README.md#create-a-bot) if you don't have.

3. Run the bot:

- Using go

```bash
go run main.go -b <your bot token> -t <notfication time (default to 17:00)> -z <timezone (default Africa/Cairo)>
```

- Using Docker

```bash
docker build -t statusbot .
docker run -e STATUSBOT_TOKEN=<your bot token> -e NOTFICATION_TIME=<notfication time> -e TIMEZONE=<timezone>-it statusbot
```

## Create a bot

- Open telegram app
- Create a new bot

```ordered
1. Find telegram bot named "@botfarther"
2. Type /newbot
```

- Get the bot token

```ordered
1. In the same bot named "@botfarther"
2. Type /token
3. Choose your bot
```
