package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	tele "gopkg.in/telebot.v4"
)

type Reminder struct {
	ID      int
	ChatID  int64
	Time    string
	Message string
}

var (
	reminders  []Reminder
	scheduler  = gocron.NewScheduler(time.Local)
	nextID     = 1
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

	// /notify HH:MM сообщение
	bot.Handle("/notify", func(c tele.Context) error {
		args := c.Args()
		if len(args) < 2 {
			return c.Send("👇 Укажи время и сообщение. Пример: /notify 14:30 покормить кота")
		}

		timeStr := args[0]
		message := strings.Join(args[1:], " ")

		if _, err := time.Parse("15:04", timeStr); err != nil {
			return c.Send("⛔ Время должно быть в формате HH:MM, например: 14:30")
		}

		chatID := c.Chat().ID

		reminder := Reminder{
			ID:      nextID,
			ChatID:  chatID,
			Time:    timeStr,
			Message: message,
		}
		reminders = append(reminders, reminder)
		nextID++

		scheduler.Every(1).Day().At(timeStr).Do(func(chatID int64, message string) {
			bot.Send(tele.ChatID(chatID), "❗ Напоминание: "+message)
		}, chatID, message)
		scheduler.StartAsync()

		return c.Send(fmt.Sprintf("✅ Напоминание #%d установлено на %s: %s", reminder.ID, timeStr, message))
	})

	// /list — показать напоминания пользователя
	bot.Handle("/list", func(c tele.Context) error {
		chatID := c.Chat().ID
		var result []string

		for _, r := range reminders {
			if r.ChatID == chatID {
				result = append(result, fmt.Sprintf("[%d] 🕒 %s — %s", r.ID, r.Time, r.Message))
			}
		}

		if len(result) == 0 {
			return c.Send("❌ У тебя нет активных напоминаний.")
		}

		return c.Send("📋 Твои напоминания:\n" + strings.Join(result, "\n"))
	})

	// /delete ID — удалить напоминание по ID
	bot.Handle("/delete", func(c tele.Context) error {
		args := c.Args()
		if len(args) < 1 {
			return c.Send("⛔ Укажи ID напоминания. Пример: /delete 1")
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return c.Send("❗ ID должен быть числом. Пример: /delete 1")
		}

		chatID := c.Chat().ID
		var updated []Reminder
		found := false

		for _, r := range reminders {
			if r.ID == id && r.ChatID == chatID {
				found = true
				continue // пропускаем для удаления
			}
			updated = append(updated, r)
		}

		if !found {
			return c.Send("❌ Напоминание с таким ID не найдено.")
		}

		reminders = updated

		// Перезапуск задач
		scheduler.Clear()
		for _, r := range reminders {
			scheduler.Every(1).Day().At(r.Time).Do(func(chatID int64, msg string) {
				bot.Send(tele.ChatID(chatID), "❗ Напоминание: "+msg)
			}, r.ChatID, r.Message)
		}
		scheduler.StartAsync()

		return c.Send(fmt.Sprintf("🗑 Напоминание #%d удалено.", id))
	})

	// /clear — удалить ВСЕ напоминания пользователя
	bot.Handle("/clear", func(c tele.Context) error {
		chatID := c.Chat().ID
		var remaining []Reminder

		for _, r := range reminders {
			if r.ChatID != chatID {
				remaining = append(remaining, r)
			}
		}

		reminders = remaining

		// Перезапуск задач
		scheduler.Clear()
		for _, r := range reminders {
			scheduler.Every(1).Day().At(r.Time).Do(func(chatID int64, msg string) {
				bot.Send(tele.ChatID(chatID), "❗ Напоминание: "+msg)
			}, r.ChatID, r.Message)
		}
		scheduler.StartAsync()

		return c.Send("🧹 Все твои напоминания очищены.")
	})

	// /help — показать список команд
	bot.Handle("/help", func(c tele.Context) error {
		return c.Send(`📖 Список команд:

/notify HH:MM сообщение — добавить напоминание
/list — показать список напоминаний
/delete ID — удалить напоминание по ID
/clear — удалить все напоминания
/help — справка`)
	})

	log.Println("Бот запущен...")
	bot.Start()
}
