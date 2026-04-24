package main

import (
	"fknow/commands"
	"fknow/listeners"
	"fknow/utils"
	"os"
	"time"

	log "github.com/charmbracelet/log"
	"github.com/eduardolat/goeasyi18n"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file, reading from environment")
	}

	if err := utils.InitTokens(os.Getenv("KNOWUNITY_REFRESH_TOKEN")); err != nil {
		log.Fatal("Knowunity token init failed", "err", err)
	}

	i18n := utils.NewI18nWithEscape(goeasyi18n.NewI18n())

	itTranslations, err := goeasyi18n.LoadFromJsonFiles("./locale/it/messages.json")
	if err != nil {
		log.Fatal("Failed to load translations", "err", err)
	}
	i18n.AddLanguage("it", itTranslations)

	b, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("Token"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("Failed to start bot", "err", err)
	}

	b.Handle("/start", func(c tele.Context) error { return commands.Start(c, i18n) })
	b.Handle("/help", func(c tele.Context) error { return commands.Help(c, i18n) })
	b.Handle(tele.OnText, func(c tele.Context) error { return listeners.OnText(c, i18n) })
	b.Handle(tele.OnAddedToGroup, func(c tele.Context) error { return listeners.OnAddedToGroup(c, i18n) })
	b.Handle(tele.OnQuery, func(c tele.Context) error { return listeners.OnQuery(c, i18n) })

	log.Info("Bot started")
	b.Start()
}
