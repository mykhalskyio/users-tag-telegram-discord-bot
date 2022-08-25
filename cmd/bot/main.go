package main

import (
	"log"

	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/bridge"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/config"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/discord"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/storage"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/telegram"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/pkg/queue"
)

func main() {
	cfg, help, err := config.GetConfig("config.yml")
	if err != nil {
		log.Fatalln(err, help)
	}
	db, err := storage.NewPostgres(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	queue, err := queue.GetQueue(cfg.Kafka.Address, cfg.Kafka.Topic)
	if err != nil {
		log.Fatalln(err)
	}
	dbot, err := discord.NewDiscordBot(cfg.Discord.Token, cfg.Discord.Prefix, queue)
	if err != nil {
		log.Fatalln(err)
	}
	tbot, err := telegram.NewTelegramBot(cfg.Telegram.Token, db, queue)
	if err != nil {
		log.Fatalln(err)
	}
	bridge := bridge.NewBridge(tbot.Api, dbot.Api, queue)

	go tbot.Start()
	go dbot.Start()
	go bridge.Start(cfg.Kafka.Address, cfg.Kafka.Topic)

	<-make(chan struct{})
}
