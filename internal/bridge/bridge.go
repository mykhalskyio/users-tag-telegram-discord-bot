package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/entity"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/pkg/queue"
	"github.com/segmentio/kafka-go"
)

type Bridge struct {
	TBot  *tgbotapi.BotAPI
	DBot  *discordgo.Session
	Queue *queue.Queue
}

func NewBridge(tBot *tgbotapi.BotAPI, dBot *discordgo.Session, queue *queue.Queue) *Bridge {
	return &Bridge{
		TBot:  tBot,
		DBot:  dBot,
		Queue: queue,
	}
}

func (b *Bridge) Start(address string, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{address},
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
		}
		fmt.Println(msg.Value)
	}
}

func (b *Bridge) SendDiscordQueue(channelId string, text string) {
	msg := entity.Message{
		Text:               text,
		Discord_channel_id: channelId,
	}
	msgJson, _ := json.Marshal(msg)
	b.Queue.Kafka.WriteMessages(
		kafka.Message{Value: msgJson},
	)
}

func (b *Bridge) SendTelegramQueue(chatId int64, text string) {
	msg := entity.Message{
		Text:             text,
		Telegram_chat_id: chatId,
	}
	msgJson, _ := json.Marshal(msg)
	b.Queue.Kafka.WriteMessages(
		kafka.Message{Value: msgJson},
	)
}
