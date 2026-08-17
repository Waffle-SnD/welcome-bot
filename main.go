package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	BotToken         string `json:"bot_token"`
	WelcomeChannelID string `json:"welcome_channel_id"`
	UserRoleID       string `json:"user_role_id"`
}

type Bot struct {
	session *discordgo.Session
	config  *Config
}

func load() (*Config, error) {
	data, err := os.ReadFile("dc_config.json")
	if err != nil {
		return nil, fmt.Errorf("read config: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %v", err)
	}

	if cfg.BotToken == "" {
		return nil, fmt.Errorf("bot token is missing")
	}

	if cfg.WelcomeChannelID == "" {
		return nil, fmt.Errorf("welcome channel id is missing")
	}

	if cfg.UserRoleID == "" {
		return nil, fmt.Errorf("user role id is missing")
	}

	return &cfg, nil
}

func (b *Bot) onGuildMemberAdd(
	s *discordgo.Session,
	member *discordgo.GuildMemberAdd,
) {
	err := s.GuildMemberRoleAdd(
		member.GuildID, 
		member.User.ID, 
		b.config.UserRoleID,
	)
	if err != nil {
		log.Printf("failed to give role to %s (%s): %v", member.User.Username, member.User.ID, err)
		return
	}

	_, err = s.ChannelMessageSend(
		b.config.WelcomeChannelID,
		fmt.Sprintf("Welcome to the server, <@%s>! 👋", member.User.ID),
	)

	if err != nil {
		log.Printf(
			"failed to send welcome message for %s (%s): %v",
			member.User.Username,
			member.User.ID,
			err,
		)
		return
	}
}

func (b *Bot) stop() {
	if err := b.session.Close(); err != nil {
		log.Printf("failed to close Discord session: %v", err)
	}
}

func main() {
	cfg, err := load()
	if err != nil {
		log.Fatal(err)
	}

	sess, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Fatal("create Discord session: %v", err)
	}

	bot := &Bot{session: sess, config: cfg}
	bot.session.AddHandler(bot.onGuildMemberAdd)
	bot.session.Identify.Intents = discordgo.IntentsGuildMembers

	if err := bot.session.Open(); err != nil {
		log.Fatal("open Discord connection: %v", err)
	}
	defer bot.stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop
}
