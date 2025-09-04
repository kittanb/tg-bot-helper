package main

import (
    "github.com/go-co-op/gocron"
    tele "gopkg.in/telebot.v4"
    "time"
)

type ReminderScheduler struct {
    scheduler *gocron.Scheduler
    bot       *tele.Bot
}

func NewReminderScheduler(bot *tele.Bot) *ReminderScheduler {
    return &ReminderScheduler{
        scheduler: gocron.NewScheduler(time.Local),
        bot:       bot,
    }
}

func (rs *ReminderScheduler) Schedule(reminders []Reminder) {
    rs.scheduler.Clear()
    for _, r := range reminders {
        rs.scheduler.Every(1).Day().At(r.Time).Do(func(chatID int64, msg string) {
            rs.bot.Send(tele.ChatID(chatID), "❗ Напоминание: "+msg)
        }, r.ChatID, r.Message)
    }
    rs.scheduler.StartAsync()
}
