package domain

import "time"

type User struct {
ID       string    `json:"id"`
Name     string    `json:"name"`
Email    string    `json:"email"`
Role     string    `json:"role"`
Password string    `json:"-"`
IsAdmin  bool      `json:"isAdmin"`
Date     time.Time `json:"date"`
}

type Policy struct {
ID        string    `json:"id"`
Role      string    `json:"role"`
BanList   []string  `json:"banList"`
CreatedAt time.Time `json:"createdAt"`
}

type FavoriteProcess struct {
ID         string    `json:"id"`
ProcessList []string `json:"processList"`
CreatedAt  time.Time `json:"createdAt"`
}

type LogEntry struct {
Timestamp string         `json:"timestamp"`
Level     string         `json:"level"`
Message   string         `json:"message"`
Meta      map[string]any `json:"meta,omitempty"`
}
