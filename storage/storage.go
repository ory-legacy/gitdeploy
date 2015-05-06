package storage

import "time"

type Storage interface {
	AddApp(app string, expiresAt time.Time, repository string) (*App, error)
	GetApp(id string) (*App, error)
	GetAppKillList() ([]*App, error)
	KillApp(app *App) error
	AddDeployEvent(app, message string) (*DeployEvent, error)
	GetNextUnreadMessage(app string) (*DeployEvent, error)
	DeployEventIsRead(event *DeployEvent) error
}

type App struct {
	ID         string    `json:"id",bson:"id"`
	ExpiresAt  time.Time `json:"expiresAt",bson:"expiresat"`
	CreatedAt  time.Time `json:"createdAt",bson:"createdat"`
	URL        string    `json:"url",bson:"url"`
	Repository string    `json:"repository",bson:"repository"`
	Killed     bool      `json:"killed",bson:"killed"`
}

type DeployEvent struct {
	ID        string    `json:"id",bson:"id"`
	App       string    `json:"app",bson:"app"`
	Message   string    `json:"message",bson:"message"`
	Timestamp time.Time `json:"timestamp",bson:"timestamp"`
	Unread    bool      `json:"unread",bson:"unread"`
}
