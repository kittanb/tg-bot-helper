package main

import (
    "sync"
)

type Reminder struct {
    ID      int
    ChatID  int64
    Time    string
    Message string
}

type ReminderStore struct {
    mu        sync.Mutex
    reminders []Reminder
    nextID    int
}

func NewReminderStore() *ReminderStore {
    return &ReminderStore{nextID: 1}
}

func (s *ReminderStore) Add(r Reminder) Reminder {
    s.mu.Lock()
    defer s.mu.Unlock()
    r.ID = s.nextID
    s.nextID++
    s.reminders = append(s.reminders, r)
    return r
}

func (s *ReminderStore) Delete(chatID int64, id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    found := false
    var updated []Reminder
    for _, r := range s.reminders {
        if r.ID == id && r.ChatID == chatID {
            found = true
            continue
        }
        updated = append(updated, r)
    }
    s.reminders = updated
    return found
}

func (s *ReminderStore) List(chatID int64) []Reminder {
    s.mu.Lock()
    defer s.mu.Unlock()
    var result []Reminder
    for _, r := range s.reminders {
        if r.ChatID == chatID {
            result = append(result, r)
        }
    }
    return result
}

func (s *ReminderStore) Clear(chatID int64) {
    s.mu.Lock()
    defer s.mu.Unlock()
    var remaining []Reminder
    for _, r := range s.reminders {
        if r.ChatID != chatID {
            remaining = append(remaining, r)
        }
    }
    s.reminders = remaining
}

func (s *ReminderStore) All() []Reminder {
    s.mu.Lock()
    defer s.mu.Unlock()
    return append([]Reminder(nil), s.reminders...)
}
