package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	gorillasession "github.com/gorilla/sessions"
	"github.com/ory-am/common/env"
	"github.com/ory-am/common/mgopath"
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/eco"
	"github.com/ory-am/gitdeploy/ip"
	"github.com/ory-am/gitdeploy/job"
	"github.com/ory-am/gitdeploy/job/deploy"
	gdLog "github.com/ory-am/gitdeploy/log"
	"github.com/ory-am/gitdeploy/public"
	"github.com/ory-am/gitdeploy/sse"
	"github.com/ory-am/gitdeploy/storage"
	"github.com/ory-am/gitdeploy/storage/mongo"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/google-json-style-response/responder"
	"gopkg.in/mgo.v2"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"time"
)

const (
	sessionCurrentDeployment = "cdid"
	sessionName              = "gdp"
)

var (
	ApiVersion      = "1.0"
	host            = env.Getenv("HOST", "")
	port            = env.Getenv("PORT", "7654")
	envAppTtl       = env.Getenv("APP_TTL", "30m")
	envClusterConf  = env.Getenv("FLYNN_CLUSTER_CONFIG", "")
	corsAllowOrigin = env.Getenv("CORS_ALLOW_ORIGIN", "http://localhost:9000")
	sessionStore    = gorillasession.NewCookieStore([]byte(env.Getenv("SESSION_SECRET", "changme")))
	envMongoPath    = env.Getenv("MONGODB", "mongodb://localhost:27017/gitdeploy")
)

type deployRequest struct {
	Repository string `json:"repository",validate:"regex=^(https|http):\\/\\/github\\.com\\/[a-zA-Z0-9\\-\\_\\.]+/[a-zA-Z0-9\\-\\_\\.]+$"`
	Ref        string `json:"ref",validate:"regex=^(origin\\/.*|tags\\/.*|[a-z0-9]+)$"`
}

type appResponse struct {
	App string `json:"app"`
}

func main() {
	eco.IsGitAvailable()
	eco.IsFlynnAvailable()
	if b := flag.Bool("init", false, "Initialize flynn and git"); *b {
		eco.InitGit()
		eco.InitFlynn(envClusterConf)
	}

	db, dbName, err := mgopath.Connect(envMongoPath)
	if err != nil {
		log.Fatal(err)
	}
	eventManager := event.New()
	storage := mongo.New(db, dbName)
	sseBroker := sse.New(storage)
	eventManager.AttachListenerAggregate(storage)
	eventManager.AttachListenerAggregate(new(gdLog.Listener))

	// Mux router
	r := mux.NewRouter()
	r.HandleFunc("/config", configHandler).Methods("GET")
	r.HandleFunc("/deployments", deployWrapperAction(sseBroker, eventManager, storage)).Methods("POST")
	r.HandleFunc("/deployments", setCORSHeaders).Methods("OPTIONS")
	r.HandleFunc("/deployments/{app:.+}/events", eventWrapperAction(sseBroker)).Methods("GET")
	r.HandleFunc("/apps/{app:.+}", getAppHandler(storage)).Methods("GET")
	r.PathPrefix("/").HandlerFunc(public.HTML5ModeHandler("./app/dist", "index.html"))
	http.Handle("/", r)

	go job.KillAppsOnHitList(storage)

	listen := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listening on %s", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}

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

	dr := new(deployRequest)
	app := uuid.NewRandom().String()
	ttl, err := time.ParseDuration(envAppTtl)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	} else if err := json.NewDecoder(r.Body).Decode(dr); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	} else if err := validator.Validate(dr); err != nil {
		responseError(w, http.StatusBadRequest, fmt.Sprintf("Validation failed: %s", err))
		return
	}

	dr.Repository = dr.Repository + ".git"
	appEntity, err := store.AddApp(app, time.Now().Add(ttl), dr.Repository, ip.GetRemoteAddr(r), dr.Ref)
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

func cleanUpSession(w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Please delete your cookies.")
		return err
	}
	delete(session.Values, sessionCurrentDeployment)
	if err := session.Save(r, w); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return err
	}
	return nil
}

// Set the different CORS headers required for CORS request
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

func responseError(w http.ResponseWriter, code int, message string) {
	log.Printf("Error %d: %s", code, message)
	response := responder.New(ApiVersion)
	response.Write(w, response.Error(code, message))
}

func responseSuccess(w http.ResponseWriter, data interface{}) {
	response := responder.New(ApiVersion)
	response.Write(w, response.Success(data))
}
