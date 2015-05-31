// Package storage abstracts the database API
package storage

import "time"

// Storage is GitDeploy's persistent API
type Storage interface {
	// AddApp adds an app to the database. Returns the app and errors from the database.
	AddApp(app string, expiresAt time.Time, repository, ip, ref string) (*App, error)
	GetApp(id string) (*App, error)
	FindAppsOnKillList() ([]*App, error)
	KillApp(app *App) error
	AddDeployEvent(app, message string) (*DeployEvent, error)
	FindDeployLogsForApp(app string) ([]*DeployEvent, error)
	GetNextUnreadDeployEvent(app string) (*DeployEvent, error)
	DeployEventIsRead(event *DeployEvent) error
	AddAppliance(app, appliance, name string) (*Appliance, error)
}

// App is the app entity.
type App struct {
	ID         string      `json:"id",bson:"id"`
	ExpiresAt  time.Time   `json:"expiresAt",bson:"expiresat"`
	CreatedAt  time.Time   `json:"createdAt",bson:"createdat"`
	URL        string      `json:"url",bson:"url"`
	Repository string      `json:"repository",bson:"repository"`
	Killed     bool        `json:"killed",bson:"killed"`
	IP         string      `json:"ip",bson:"ip"`
	Appliances []Appliance `json:"appliances",bson:"appliances"`
	Ref string `json:"ref",bson:"ref"`
}

// Appliance is the appliance entity.
type Appliance struct {
	ID        string    `json:"id",bson:"id"`
	Name      string    `json:"name",bson:"name"`
	Killed    bool      `json:"killed",bson:"killed"`
	CreatedAt time.Time `json:"createdAt",bson:"createdat"`
}

// DeployEvent saves the deployment logs/events.
type DeployEvent struct {
	ID        string    `json:"id",bson:"id"`
	App       string    `json:"app",bson:"app"`
	Message   string    `json:"message",bson:"message"`
	Timestamp time.Time `json:"timestamp",bson:"timestamp"`
	Unread    bool      `json:"unread",bson:"unread"`
}
