package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/ip"
	"github.com/ory-am/gitdeploy/job/deploy"
	"github.com/ory-am/gitdeploy/sse"
	"github.com/ory-am/gitdeploy/storage"
	"github.com/ory-am/gitdeploy/storage/mongo"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/flynn"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"regexp"
	"time"
)

func eventWrapperAction(sseBroker *sse.Broker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		sseBroker.EventHandler(w, r)
	}
}

func deployWrapperAction(sseBroker *sse.Broker, em *event.EventManager, store *mongo.MongoStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		deployAction(w, r, sseBroker, em, store)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)
	responseSuccess(w, struct {
		Now time.Time `json:"now"`
	}{time.Now()})
}

func getAppHandler(store *mongo.MongoStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if cleanUpSession(w, r) != nil {
			return
		}
		setCORSHeaders(w, r)
		vars := mux.Vars(r)
		id := vars["app"]
		app, err := store.GetApp(id)
		if err == mgo.ErrNotFound {
			responseError(w, http.StatusNotFound, fmt.Sprintf("App %s does not exist.", id))
			return
		} else if app.Killed {
			responseError(w, http.StatusNotFound, "App timed out and is no longer available.")
			return
		} else if err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		logs, err := flynn.GetLogs(app.ID)
		if err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		ps, err := flynn.GetProcs(app.ID)
		if err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		deployLogs, err := store.FindDeployLogsForApp(id)
		if err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		responseSuccess(w, struct {
			*storage.App
			Logs       string                 `json:"logs"`
			DeployLogs []*storage.DeployEvent `json:"deployLogs"`
			PS         string                 `json:"ps"`
		}{app, logs, deployLogs, ps})
	}
}

func deployAction(w http.ResponseWriter, r *http.Request, sseBroker *sse.Broker, em *event.EventManager, store *mongo.MongoStorage) {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Please delete your cookies.")
		return
	}

	// Check if the user is currently deploying an application and switch to that one.
	if v, ok := session.Values[sessionCurrentDeployment].(string); ok && len(v) > 0 {
		app, err := store.GetApp(v)
		if err != nil {
			cleanUpSession(w, r)
			log.Printf("Could not fetch app from cookie: %s", err.Error())
		} else if !sseBroker.IsChannelOpen(app.ID) {
			cleanUpSession(w, r)
			log.Printf("Channel %s does not exist any more", app.ID)
		} else {
			responseSuccess(w, app)
			return
		}
	}

	// Parse body
	dr := new(deployRequest)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(dr)

	// Validate URL
	regExpression := "(https|http):\\/\\/github\\.com\\/[a-zA-Z0-9\\-\\_\\.]+/[a-zA-Z0-9\\-\\_\\.]+\\.git"
	if match, _ := regexp.MatchString(regExpression, dr.Repository); !match {
		responseError(w, http.StatusBadRequest, "I only support GitHub.")
		return
	}

	app := uuid.NewRandom().String()
	ttl, err := time.ParseDuration(envAppTtl)
	if err != nil {
		log.Println(err.Error())
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	appEntity, err := store.AddApp(app, time.Now().Add(ttl), dr.Repository, ip.GetRemoteAddr(r))
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	session.Values[sessionCurrentDeployment] = appEntity.ID
	if err := session.Save(r, w); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseSuccess(w, appEntity)
	go func() {
		sseBroker.OpenChannel(app)
		sseBroker.Start(app)
		defer func() {
			// Give the client the chance to read the output...
			time.Sleep(15 * time.Second)
			sseBroker.CloseChannel(app)
		}()
		l := make(task.WorkerLog)
		tasks := deploy.CreateJob(l, store, appEntity)
		if err := task.RunJob(l, app, em, tasks); err != nil {
			log.Printf("RUNJOB ERROR: %s", err)
		}
	}()
}
