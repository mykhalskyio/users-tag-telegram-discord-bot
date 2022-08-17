package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/entity"
)

type TelegramStorage interface {
	Insert(username string, chatId int) error
	GetAll(chatId int) (*[]entity.ChatUser, error)
	Delete(username string, chatId int) error
	Get(username string, chatId int) (*entity.ChatUser, error)
}

type TelegramBot struct {
	Api     *tgbotapi.BotAPI
	Storage TelegramStorage
}

func NewTelegramBot(token string, storage TelegramStorage) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{
		Api:     bot,
		Storage: storage,
	}, nil
}

func (bot *TelegramBot) Start() {
	bot.Api.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.Api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if user, _ := bot.Storage.Get(update.Message.From.UserName, int(update.Message.Chat.ID)); user == nil {
				bot.Storage.Insert(update.Message.From.UserName, int(update.Message.Chat.ID))
			}
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "all_chat_users":
					users, _ := bot.Storage.GetAll(int(update.Message.Chat.ID))
					usersStr := ""
					for _, user := range *users {
						usersStr += " " + user.Username
					}
					msg.Text = usersStr
				}
				bot.Api.Send(msg)
			}
		}
	}
}
