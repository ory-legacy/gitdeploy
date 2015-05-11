package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	gorillasession "github.com/gorilla/sessions"
	"github.com/ory-am/common/env"
	"github.com/ory-am/common/mgopath"
	"github.com/ory-am/event"
	gde "github.com/ory-am/gitdeploy/event"
	"github.com/ory-am/gitdeploy/job"
	gdLog "github.com/ory-am/gitdeploy/log"
	"github.com/ory-am/gitdeploy/sse"
	"github.com/ory-am/gitdeploy/storage"
	"github.com/ory-am/gitdeploy/storage/mongo"
	"github.com/ory-am/google-json-style-response/responder"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"
	"runtime"
	"os"
    "strings"
)

const (
	sessionCurrentDeployment = "cdid"
	sessionName = "gdp"
)

var (
// API Version
	ApiVersion = "1.0"

// Generic configuration
	host = env.Getenv("HOST", "")
	port = env.Getenv("PORT", "7654")

	envAppTtl = env.Getenv("APP_TTL", "30m")

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
	r.HandleFunc("/config", configHandler).Methods("GET")
	r.HandleFunc("/deployments", deployWrapperAction(sseBroker, eventManager, storage)).Methods("POST")
	r.HandleFunc("/deployments", setCORSHeaders).Methods("OPTIONS")
	r.HandleFunc("/deployments/{app:.+}/events", eventWrapperAction(sseBroker)).Methods("GET")
	r.HandleFunc("/apps/{app:.+}", getAppHandler(storage)).Methods("GET")
	r.PathPrefix("/").HandlerFunc(publicHandler("./app/dist"))
	http.Handle("/", r)

	go job.KillAppsOnHitList(storage)

	listen := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Listening on %s", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}

func publicHandler(dir string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := dir + r.URL.Path
		if r.URL.Path == "/" {
			path = dir + "/index.html"
		}
		pattern := `!\.html|\.js|\.svg|\.css|\.png|\.jpg$`

		if f, err := os.Stat(path); err == nil {
			if !f.IsDir() {
				http.ServeFile(w, r, path)
				return
			} else {
				http.Redirect(w, r, "/", 302)
				return
			}
		}

		if matched, err := regexp.MatchString(pattern, path); err != nil {
			log.Printf("Could not exec regex: %s", err.Error())
		} else if !matched {
			http.Redirect(w, r, "/", 302)
		} else {
			http.NotFound(w, r)
		}
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)
	responseSuccess(w, struct {
		Now time.Time `json:"now"`
	}{
		time.Now(),
	})
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
		logs, err := job.GetLogs(app.ID)
		if err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}

		ps, err := job.GetPS(app.ID)
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
	if v, ok := session.Values[sessionCurrentDeployment].(string); ok && len(v) > 0 {
		app, err := store.GetApp(v)
		if err != nil {
			cleanUpSession(w, r)
			log.Printf("Could not fetch app from cookie: %s", err.Error())
		} else {
			responseSuccess(w, app)
		}
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
		appEntity, err := store.AddApp(app, time.Now().Add(ttl), dr.Repository, getIP(r))
		if err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		session.Values[sessionCurrentDeployment] = appEntity.ID
		if err := session.Save(r, w); err != nil {
			responseError(w, http.StatusInternalServerError, err.Error())
			return
		}
		go runJobs(w, r, em, dr, appEntity, sseBroker, session, store)
		responseSuccess(w, appEntity)
	}
}

func runJobs(w http.ResponseWriter, r *http.Request, em *event.EventManager, dr *deployRequest, app *storage.App, sseBroker *sse.Broker, session *gorillasession.Session, store *mongo.MongoStorage) {
	sseBroker.OpenChannel(app.ID)
	sseBroker.Start(app.ID)
	defer func() {
		// Give the client the chance to read the output...
		time.Sleep(2 * time.Minute)
		sseBroker.CloseChannel(app.ID)
	}()

	em.Trigger("app.created", gde.New(app.ID, app.ID))

	destination, err := job.Clone(em, app.ID, dr.Repository)
	if err != nil {
		log.Printf("Error in job.clone %s: %s", app.ID, err.Error())
		return
	}

	if err = job.Parse(em, app.ID, destination); err != nil {
		log.Printf("Error in job.parse %s: %s", app.ID, err.Error())
		return
	}

	if err = job.Deploy(em, app.ID, destination); err != nil {
		log.Printf("Error in job.deploy %s: %s", app.ID, err.Error())
		return
	}

	cluster, err := job.GetCluster(em, app.ID)
	if err != nil {
		log.Printf("Error in job.deploy %s: %s", app.ID, err.Error())
		return
	}

	app.URL = fmt.Sprintf("%s.%s", app.ID, cluster.Host)
	if err := store.UpdateApp(app); err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		log.Printf("Error %s: %s", app.ID, err.Error())
		return
	}

	log.Println("Deployment successful.")
	em.Trigger("app.deployed", gde.New(app.ID, app.URL))
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

func checkDependencies() {
	go func() {
		_, err := exec.LookPath("git")
		if err != nil {
			log.Fatal("Git CLI is required but not installed or not in path.")
		}
	}()
	go checkIfFlynnExists()
}

func checkIfFlynnExists() {
	func() {
		_, err := exec.LookPath("flynn")
		if err != nil {
			if runtime.GOOS == "windows" {
				log.Fatal("Flynn CLI is required but not installed or not in path.")
			}
			log.Println("Could not find Flynn CLI, trying to install...")
			if _, err := exec.Command("sh", "flynn.sh").CombinedOutput(); err != nil {
				log.Printf("Could not install Flynn CLI: %s", err.Error())
			} else if _, err := exec.LookPath("flynn"); err != nil {
				log.Fatal("Could not install Flynn CLI.")
			}
		}
	}()
}

func getIP(r *http.Request) string {
    ip := removePort(r.RemoteAddr)
    if len(r.Header.Get("X-FORWARDED-FOR")) > 0 {
        ip = r.Header.Get("X-FORWARDED-FOR")
    }
    return ip
}

func removePort(ip string) string {
    split := strings.Split(ip, ":")
    if len(split) < 2 {
        return ip
    }
    return strings.Join(split[:len(split)-1], ":")
}