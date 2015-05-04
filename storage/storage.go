package storage

import "time"

type Storage interface {
	AddApp(app string, ttl time.Time) (*App, error)
	GetApp(id string) (*App, error)
	GetAppKillList() ([]*App, error)
	KillApp(app *App) error
	AddLogEvent(app, message string) (*LogEvent, error)
	GetNextUnreadMessage(app string) (*LogEvent, error)
	LogEventIsRead(event *LogEvent) error
}

type App struct {
	ID        string    `json:"id",bson:"id"`
	TTL       time.Time `json:"ttl",bson:"ttl"`
	Timestamp time.Time `json:"ttl",bson:"ttl"`
	Killed    bool      `json:"killed",bson:"killed"`
}

type LogEvent struct {
	ID        string    `json:"id",bson:"id"`
	App       string    `json:"app",bson:"app"`
	Message   string    `json:"message",bson:"message"`
	Timestamp time.Time `json:"timestamp",bson:"timestamp"`
	Unread    bool      `json:"unread",bson:"unread"`
}
