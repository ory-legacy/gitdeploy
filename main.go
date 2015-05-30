package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	gorillasession "github.com/gorilla/sessions"
	"github.com/ory-am/common/env"
	"github.com/ory-am/common/mgopath"
	"github.com/ory-am/event"
	"github.com/ory-am/gitdeploy/eco"
	"github.com/ory-am/gitdeploy/job"
	gdLog "github.com/ory-am/gitdeploy/log"
	"github.com/ory-am/gitdeploy/public"
	"github.com/ory-am/gitdeploy/sse"
	"github.com/ory-am/gitdeploy/storage/mongo"
	"log"
	"net/http"
)

const (
	sessionCurrentDeployment = "cdid"
	sessionName              = "gdp"
)

var (
	// API Version
	ApiVersion = "1.0"

	// Generic configuration
	host = env.Getenv("HOST", "")
	port = env.Getenv("PORT", "7654")

	envAppTtl      = env.Getenv("APP_TTL", "30m")
	envClusterConf = env.Getenv("FLYNN_CLUSTER_CONFIG", "")

	// Configuration for CORS
	corsAllowOrigin = env.Getenv("CORS_ALLOW_ORIGIN", "http://localhost:9000")

	sessionStore = gorillasession.NewCookieStore([]byte(env.Getenv("SESSION_SECRET", "changme")))

	// MongoDB
	envMongoPath = env.Getenv("MONGODB", "mongodb://localhost:27017/gitdeploy")
)

type deployRequest struct {
	Repository string `json:"repository"`
}

type appResponse struct {
	App string `json:"app"`
}

func main() {
	eco.IsGitAvailable()
	eco.IsFlynnAvailable()

	if b := flag.Bool("init", false, "Initialize flynn and git"); !*b {
		eco.InitGit()
		eco.InitFlynn(envClusterConf)
	}

	eventManager := event.New()

	// mgo
	db, dbName, err := mgopath.Connect(envMongoPath)
	if err != nil {
		log.Fatal(err)
	}

	storage := mongo.New(db, dbName)
	eventManager.AttachListenerAggregate(storage)

	// SSE broker
	sseBroker := sse.New(storage)

	// Log listener
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
