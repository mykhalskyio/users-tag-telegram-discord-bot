package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Api    *discordgo.Session
	Prefix string
}

var (
	prefix string
)

func NewDiscordBot(token string, prefix string) (*DiscordBot, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &DiscordBot{
		Api:    bot,
		Prefix: prefix,
	}, nil
}

func (bot *DiscordBot) Start() {
	bot.Api.AddHandler(messageHandler)
	prefix = bot.Prefix
	bot.Api.Open()
	log.Println("Bot start")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case prefix + "ping":
		s.ChannelMessageSend(m.ChannelID, "pong!")
	case prefix + "isAdmin":
		ok, _ := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator)
		if ok {
			s.ChannelMessageSend(m.ChannelID, "Yes")
		} else {
			s.ChannelMessageSend(m.ChannelID, "False")
		}
	}
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
