package discord

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/internal/entity"
	"github.com/mykhalskyio/users-tag-telegram-discord-bot/pkg/queue"
)

type DiscordBot struct {
	Api    *discordgo.Session
	Queue  *queue.Queue
	Prefix string
}

var (
	prefix string
)

func NewDiscordBot(token string, prefix string, queue *queue.Queue) (*DiscordBot, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &DiscordBot{
		Api:    bot,
		Queue:  queue,
		Prefix: prefix,
	}, nil
}

func (bot *DiscordBot) Start() {
	bot.Api.AddHandler(bot.setActivity)
	bot.Api.AddHandler(bot.messageHandler)
	prefix = bot.Prefix
	bot.Api.Open()
}

func (bot *DiscordBot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case prefix + "ping":
		s.ChannelMessageSend(m.ChannelID, m.GuildID)
		msgJson, _ := json.Marshal(entity.Message{
			Text:             "pong!",
			Telegram_chat_id: 64,
		})
		bot.Queue.SendToQueue(msgJson)
	case prefix + "isAdmin":
		ok, _ := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator)
		if ok {
			s.ChannelMessageSend(m.ChannelID, "Yes")
		} else {
			s.ChannelMessageSend(m.ChannelID, "False")
		}
	}
}

func (bot *DiscordBot) setActivity(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.UpdateGameStatus(1, "Cписок команд - !help")
}

func memberHasPermission(s *discordgo.Session, guildID string, userID string, permission int) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if int(role.Permissions)&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}
