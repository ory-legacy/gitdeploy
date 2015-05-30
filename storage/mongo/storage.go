package mongo

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/sse"
	"github.com/ory-am/gitdeploy/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const (
	appCollection          = "app"
	appDeployLogCollection = "appEvents"
	applianceCollection    = "appAppliances"
)

type MongoStorage struct {
	session  *mgo.Session
	database string
}

// NewUserStorage creates a new database session for storing users
func New(session *mgo.Session, database string) *MongoStorage {
	s := &MongoStorage{
		session:  session,
		database: database,
	}
	s.ensureUnique(appCollection, []string{"id"})
	s.ensureUnique(appDeployLogCollection, []string{"id"})
	return s
}

func (s *MongoStorage) AddAppliance(app, appliance, name string) (*storage.Appliance, error) {
	a := &storage.Appliance{
		ID:        appliance,
		CreatedAt: time.Now(),
		Name:      name,
		Killed:    false,
	}
	return a, s.getCollection(appCollection).Update(bson.M{"id": app}, bson.M{"$push": bson.M{"appliances": a}})
}

func (s *MongoStorage) AddApp(app string, ttl time.Time, repository, ip string) (a *storage.App, err error) {
	a = &storage.App{
		ID:         app,
		ExpiresAt:  ttl,
		CreatedAt:  time.Now(),
		Killed:     false,
		Repository: repository,
		IP:         ip,
	}
	err = s.getCollection(appCollection).Insert(a)
	return a, err
}

func (s *MongoStorage) UpdateApp(app *storage.App) error {
	return s.getCollection(appCollection).Update(bson.M{"id": app.ID}, app)
}

func (s *MongoStorage) AddDeployEvent(app, message string) (*storage.DeployEvent, error) {
	e := &storage.DeployEvent{
		ID:        uuid.NewRandom().String(),
		App:       app,
		Message:   message,
		Timestamp: time.Now(),
		Unread:    true,
	}
	return e, s.getCollection(appDeployLogCollection).Insert(e)
}

func (s *MongoStorage) GetApp(id string) (app *storage.App, err error) {
	err = s.getCollection(appCollection).Find(bson.M{"id": id}).One(&app)
	return app, err
}

func (s *MongoStorage) FindDeployLogsForApp(app string) (e []*storage.DeployEvent, err error) {
	err = s.getCollection(appDeployLogCollection).Find(bson.M{"app": app}).All(&e)
	return e, err
}

func (s *MongoStorage) GetNextUnreadDeployEvent(app string) (e *storage.DeployEvent, err error) {
	e = new(storage.DeployEvent)
	err = s.getCollection(appDeployLogCollection).Find(bson.M{
		"app":    app,
		"unread": true,
	}).Sort("+timestamp").One(e)
	return e, err
}

func (s *MongoStorage) FindAppsOnKillList() (apps []*storage.App, err error) {
	err = s.getCollection(appCollection).Find(bson.M{
		"expiresat": bson.M{
			"$lt": time.Now(),
		},
		"killed": false,
	}).All(&apps)
	return apps, err
}

func (s *MongoStorage) KillApp(app *storage.App) (err error) {
	app.Killed = true
	return s.getCollection(appCollection).Update(bson.M{"id": app.ID}, app)
}

func (s *MongoStorage) DeployEventIsRead(e *storage.DeployEvent) error {
	e.Unread = false
	return s.getCollection(appDeployLogCollection).Update(bson.M{"id": e.ID}, e)
}

func (s *MongoStorage) Trigger(name string, data interface{}) {
	if e, ok := data.(*sse.Event); ok {
		if _, err := s.AddDeployEvent(e.App, e.SSEify()); err != nil {
			log.Printf(err.Error())
		}
	}
}

func (l *MongoStorage) AttachAggregate(em *event.EventManager) {
	em.AttachListener("git.clone", l)
	em.AttachListener("config.parse", l)
	em.AttachListener("config.procs", l)
	em.AttachListener("config.buildpack", l)
	em.AttachListener("config.env", l)
	em.AttachListener("config.appliances", l)
	em.AttachListener("env.commit", l)
	em.AttachListener("git.add", l)
	em.AttachListener("git.commit", l)
	em.AttachListener("app.release", l)
	em.AttachListener("app.cleanup", l)
	em.AttachListener("app.deployed", l)
}

func (s *MongoStorage) getCollection(name string) *mgo.Collection {
	return s.session.DB(s.database).C(name)
}

func (s *MongoStorage) ensureUnique(collection string, keys []string) {
	c := s.getCollection(collection)
	if err := c.EnsureIndex(mgo.Index{Key: keys, Unique: true}); err != nil {
		log.Fatalf("Could not ensure index: %s", err)
	}
}
