package telegram

import (
	"encoding/json"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/entity"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/pkg/queue"
)

type TelegramStorage interface {
	InsertUser(username string, chatId int) error
	GetAllUsers(chatId int) (*[]entity.ChatUser, bool, error)
	DeleteUser(username string, chatId int) error
	GetUser(username string, chatId int) (*entity.ChatUser, bool, error)
}

type TelegramBot struct {
	Api     *tgbotapi.BotAPI
	Queue   *queue.Queue
	Storage TelegramStorage
}

func NewTelegramBot(token string, storage TelegramStorage, queue *queue.Queue) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{
		Api:     bot,
		Queue:   queue,
		Storage: storage,
	}, nil
}

func (bot *TelegramBot) Start() {
	bot.Api.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.Api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "all_chat_users":
					users, ok, _ := bot.Storage.GetAllUsers(int(update.Message.Chat.ID))
					usersStr := ""
					if ok {
						for _, user := range *users {
							usersStr += " " + user.Username
						}
					}
					msg.Text = usersStr
				case "ping":
					msg.Text = "pong!"
					msgJson, _ := json.Marshal(entity.Message{
						Text:               "pong!",
						Discord_channel_id: "channelID",
					})
					bot.Queue.SendToQueue(msgJson)
				}
				bot.Api.Send(msg)
			} else {
				if update.Message.Text == "" {
					_, ok, _ := bot.Storage.GetUser(update.Message.From.UserName, int(update.Message.Chat.ID))
					if ok {
						log.Println("Delete " + update.Message.LeftChatMember.UserName)
						bot.Storage.DeleteUser(update.Message.LeftChatMember.UserName, int(update.Message.Chat.ID))
					}
				} else {
					if _, ok, _ := bot.Storage.GetUser(update.Message.From.UserName, int(update.Message.Chat.ID)); !ok {
						bot.Storage.InsertUser(update.Message.From.UserName, int(update.Message.Chat.ID))
					}
				}
			}
		}
	}
}
