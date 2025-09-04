package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v4"
)

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	store := NewReminderStore()
	scheduler := NewReminderScheduler(bot)

	RegisterHandlers(bot, store, scheduler)

	log.Println("Бот запущен...")
	bot.Start()
}
