package mongo

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"github.com/ory-am/gitdeploy/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
	"fmt"
)

const (
	appCollection         = "app"
	appDeployLogCollection = "appEvents"
)

type MongoStorage struct {
	db *mgo.Database
}

// NewUserStorage creates a new database session for storing users
func New(database *mgo.Database) *MongoStorage {
	ensureIndex(database.C(appCollection), mgo.Index{Key: []string{"id"}, Unique: true})
	ensureIndex(database.C(appDeployLogCollection), mgo.Index{Key: []string{"id"}, Unique: true})
	return &MongoStorage{database}
}

func (s *MongoStorage) AddApp(app string, ttl time.Time, repository string) (a *storage.App, err error) {
	a = &storage.App{
		ID:         app,
		ExpiresAt:  ttl,
		CreatedAt:  time.Now(),
		Killed:     false,
		Repository: repository,
	}
	err = s.db.C(appCollection).Insert(a)
	return a, err
}

func (s *MongoStorage) UpdateApp(app *storage.App) error {
	return s.db.C(appCollection).Update(bson.M{"id": app.ID}, app)
}

func (s *MongoStorage) AddDeployEvent(app, message string) (*storage.DeployEvent, error) {
	e := &storage.DeployEvent{
		ID:        uuid.NewRandom().String(),
		App:       app,
		Message:   message,
		Timestamp: time.Now(),
		Unread:    true,
	}
	return e, s.db.C(appDeployLogCollection).Insert(e)
}

func (s *MongoStorage) GetApp(id string) (app *storage.App, err error) {
	err = s.db.C(appCollection).Find(bson.M{"id": id}).One(&app)
	return app, err
}

func (s *MongoStorage) FindDeployLogsForApp(app string) (e []*storage.DeployEvent, err error) {
	err = s.db.C(appDeployLogCollection).Find(bson.M{"app": app}).All(&e)
	return e, err
}

func (s *MongoStorage) GetNextUnreadDeployEvent(app string) (e *storage.DeployEvent, err error) {
	e = new(storage.DeployEvent)
	err = s.db.C(appDeployLogCollection).Find(bson.M{
		"app":    app,
		"unread": true,
	}).Sort("+timestamp").One(e)
	return e, err
}

func (s *MongoStorage) FindAppsOnKillList() (apps []*storage.App, err error) {
	err = s.db.C(appCollection).Find(bson.M{
		"expiresat": bson.M{
			"$lt": time.Now(),
		},
		"killed": false,
	}).All(&apps)
	return apps, err
}

func (s *MongoStorage) KillApp(app *storage.App) (err error) {
	app.Killed = true
	return s.db.C(appCollection).Update(bson.M{"id": app.ID}, app)
}

func (s *MongoStorage) DeployEventIsRead(e *storage.DeployEvent) error {
	e.Unread = false
	return s.db.C(appDeployLogCollection).Update(bson.M{"id": e.ID}, e)
}

func (s *MongoStorage) Trigger(name string, data interface{}) {
	if e, ok := data.(gde.JobEvent); ok {
		// TODO Ugly...
		e.SetEventName(name)
		if _, err := s.AddDeployEvent(e.GetApp(), e.GetMessage()); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func (s *MongoStorage) AttachAggregate(em *event.EventManager) {
	em.AttachListener("jobs.clone", s)
	em.AttachListener("jobs.deploy", s)
	em.AttachListener("jobs.parse", s)
	em.AttachListener("app.created", s)
	em.AttachListener("app.deployed", s)
	em.AttachListener("jobs.cluster", s)
}

func ensureIndex(c *mgo.Collection, i mgo.Index) {
	if err := c.EnsureIndex(i); err != nil {
		log.Fatalf("Could not ensure index: %s", err)
	}
}
