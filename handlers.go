package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
)

func RegisterHandlers(bot *tele.Bot, store *ReminderStore, scheduler *ReminderScheduler) {
	// /notify HH:MM сообщение
	bot.Handle("/notify", func(c tele.Context) error {
		args := c.Args()
		if len(args) < 2 {
			return c.Send("👇 Укажи время и сообщение. Пример: /notify 21:00 сосал таблетку?")
		}

		timeStr := args[0]
		message := strings.Join(args[1:], " ")

		if _, err := time.Parse("15:04", timeStr); err != nil {
			return c.Send("⛔ Время должно быть в формате HH:MM, например: 14:30")
		}

		reminder := Reminder{
			ChatID:  c.Chat().ID,
			Time:    timeStr,
			Message: message,
		}
		r := store.Add(reminder)
		scheduler.Schedule(store.All())

		return c.Send(fmt.Sprintf("✅ Напоминание #%d установлено на %s: %s", r.ID, timeStr, message))
	})

	// /list — показать напоминания пользователя
	bot.Handle("/list", func(c tele.Context) error {
		chatID := c.Chat().ID
		var result []string

		for _, r := range store.List(chatID) {
			result = append(result, fmt.Sprintf("[%d] 🕒 %s — %s", r.ID, r.Time, r.Message))
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
		if !store.Delete(chatID, id) {
			return c.Send("❌ Напоминание с таким ID не найдено.")
		}

		scheduler.Schedule(store.All())
		return c.Send(fmt.Sprintf("🗑 Напоминание #%d удалено.", id))
	})

	// /clear — удалить ВСЕ напоминания пользователя
	bot.Handle("/clear", func(c tele.Context) error {
		chatID := c.Chat().ID
		store.Clear(chatID)
		scheduler.Schedule(store.All())
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
}
