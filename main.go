package main

import (
	"fknow/commands"
	"fknow/listeners"
	"fknow/utils"
	"os"
	"time"

	log "github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	"github.com/eduardolat/goeasyi18n"
	tele "gopkg.in/telebot.v4"
)

func main() {
	godotenv.Load()
	i18n := utils.NewI18nWithEscape(goeasyi18n.NewI18n())

	// loading bot
	pref := tele.Settings{
		Token:  os.Getenv("Token"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// loading messages
	itTranslations, err := goeasyi18n.LoadFromJsonFiles("./locale/it/messages.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	i18n.AddLanguage("it", itTranslations)

	// commands
	b.Handle("/start", func(c tele.Context) error {
		return commands.Start(c, i18n)
	})
	b.Handle("/help", func(c tele.Context) error {
		return commands.Help(c, i18n)
	})

	// listeners
	b.Handle(tele.OnText, func(c tele.Context) error {
		return listeners.OnText(c, i18n)
	})
	b.Handle(tele.OnAddedToGroup, func(c tele.Context) error {
		return listeners.OnAddedToGroup(c, i18n)
	})
	b.Handle(tele.OnQuery, func(c tele.Context) error {
		return listeners.OnQuery(c, i18n)
	})

	b.Start()
}
