package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ory-am/common/env"
	"github.com/ory-am/common/mgopath"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"github.com/ory-am/gitdeploy/job"
	gdLog "github.com/ory-am/gitdeploy/log"
	"github.com/ory-am/gitdeploy/sse"
	"github.com/ory-am/gitdeploy/storage/mongo"
	"github.com/ory-am/google-json-style-response/responder"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"
	gorillasession "github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"
	"github.com/ory-am/gitdeploy/storage"
)

var (
	// API Version
	ApiVersion = "1.0"

	// Generic configuration
	host = env.Getenv("HOST", "")
	port = env.Getenv("PORT", "7654")

	envAppTtl = env.Getenv("APP_TTL", "5m")

	// Configuration for CORS
	corsAllowOrigin = env.Getenv("CORS_ALLOW_ORIGIN", "http://localhost:9000")

	sessionStore = gorillasession.NewCookieStore([]byte(env.Getenv("SESSION_SECRET", "changme")))
	sessionName = "gitdeploy"

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
	checkDependencies()
	eventManager := event.New()

	// mgo
	mongoDatabase, err := mgopath.Connect(envMongoPath)
	if err != nil {
		log.Fatal(err)
	}

	storage := mongo.New(mongoDatabase)
	eventManager.AttachListenerAggregate(storage)

	// SSE broker
	sseBroker := sse.New(storage)

	// Log listener
	eventManager.AttachListenerAggregate(new(gdLog.Listener))

	// Mux router
	r := mux.NewRouter()
	r.HandleFunc("/deployments", deployWrapperAction(sseBroker, eventManager, storage)).Methods("POST")
	r.HandleFunc("/deployments", setCORSHeaders).Methods("OPTIONS")
	r.HandleFunc("/deployments/{app:.+}/events", eventWrapperAction(sseBroker)).Methods("GET")
	r.HandleFunc("/apps/{app:.+}", getAppHandler(sseBroker)).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./app/")))
	http.Handle("/", r)

	go job.KillAppsOnHitList(storage)

	listen := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listening on %s", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}

func getAppHandler(store *mongo.MongoStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		vars := mux.Vars(r)
		id := vars["app"]
		app, err := store.GetApp(id)
		if err == mgo.ErrNotFound {
			responseError(w, http.StatusNotFound, "App could not be found.")
			return
		} else if err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		logs, err := job.GetLogs(app)
		if err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		responseSuccess(w, struct{
			*storage.App
			Logs string `json:"logs"`
		}{
			app,
			logs,
		})
	}
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

func deployAction(w http.ResponseWriter, r *http.Request, sseBroker *sse.Broker, em *event.EventManager, store *mongo.MongoStorage) {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		responseError(w, http.StatusBadRequest, "Please delete your cookies.")
		return
	}

	// Check if the user is currently deploying an application and switch to that one.
	if v, ok := session.Values["currentDeploymentID"]; ok {
		app, err := store.GetApp(v)
		if err != nil {
			responseError(w, http.StatusInternalServerError, "Please delete your cookies.")
			return
		}
		responseSuccess(w, app)
	}

	// Parse body
	dr := new(deployRequest)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(dr)

	// Validate URL
	regExpression := "(https|http):\\/\\/github\\.com\\/[a-zA-Z0-9\\-\\_]+/[a-zA-Z0-9\\-\\_]+\\.git"
	if match, _ := regexp.MatchString(regExpression, dr.Repository); !match {
		responseError(w, http.StatusBadRequest, "I only support GitHub.")
		return
	}

	app := uuid.NewRandom().String()
	if ttl, err := time.ParseDuration(envAppTtl); err != nil {
		log.Println(err.Error())
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		appEntity, err := store.AddApp(app, time.Now().Add(ttl))
		if err != nil {
			log.Println(err.Error())
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		session.Values["currentDeploymentID"] = appEntity.ID
		if err := session.Save(r, w); err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		responseSuccess(w, appEntity)
	}

	go runJobs(w, r, em, dr, app, sseBroker, session)
}

func runJobs(w http.ResponseWriter, r *http.Request, em *event.EventManager, dr *deployRequest, app string, sseBroker *sse.Broker, session gorillasession.Session) {
	sseBroker.OpenChannel(app)
	sseBroker.Start(app)
	defer func() {
		// Give the client the chance to read the output...
		time.Sleep(2 * time.Minute)
		delete(session.Values, "currentDeploymentID")
		if err := session.Save(r, w); err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		sseBroker.CloseChannel(app)
	}()

	em.Trigger("app.created", gde.New(app, app))

	destination, err := job.Clone(em, app, dr.Repository)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		log.Printf("Error in job.clone %s: %s", app, err.Error())
		return
	}

	if err = job.Parse(em, app, destination); err != nil {
		log.Printf("Error in job.parse %s: %s", app, err.Error())
		return
	}

	if err = job.Deploy(em, app, destination); err != nil {
		log.Printf("Error in job.deploy %s: %s", app, err.Error())
		return
	}

	cluster, err := job.GetCluster(em, app, destination)
	if err != nil {
		log.Printf("Error in job.deploy %s: %s", app, err.Error())
		return
	}

	log.Println("Deployment successful.")
	em.Trigger("app.deployed", gde.New(app, fmt.Sprintf("%s.%s", app, cluster)))
}

// Set the different CORS headers required for CORS request
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

func responseError(w http.ResponseWriter, code int, message string) {
	response := responder.New(ApiVersion)
	response.Write(w, response.Error(code, message))
}

func responseSuccess(w http.ResponseWriter, data interface{}) {
	response := responder.New(ApiVersion)
	response.Write(w, response.Success(data))
}

func checkDependencies() {
	go func() {
		_, err := exec.LookPath("git")
		if err != nil {
			log.Fatal("Git CLI is required but not installed or not in path.")
		}
	}()
	go func() {
		_, err := exec.LookPath("flynn")
		if err != nil {
			log.Fatal("Flynn CLI is required but not installed or not in path.")
		}
	}()
	go func() {
		_, err := exec.LookPath("bash")
		if err != nil {
			log.Fatal("Bash is required but not installed or not in path.")
		}
	}()
}
