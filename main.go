package main

import (
    "code.google.com/p/go-uuid/uuid"
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/ory-am/common/env"
    "github.com/ory-am/gitdeploy/flynn"
    "github.com/ory-am/gitdeploy/github"
    "github.com/ory-am/gitdeploy/sse"
    "github.com/ory-am/google-json-style-response/responder"
    "log"
    "net/http"
    "net/url"
    "os"
    "runtime"
    "os/exec"
    "time"
    "database/sql"
    _ "github.com/lib/pq"
    "io/ioutil"
)

const apiVersion = "1.0"

type deployRequest struct {
    Repository string `json:"repository"`
}

type pipe struct {}

func (p *pipe) Write(d []byte) (n int, err error) {
    log.Printf("Git subcommand: %s", d)
    return len(d), nil
}

func main() {
    pgHost := env.Getenv("PGHOST", "localhost")
    pgDatabase := env.Getenv("PGDATABASE", "gitdeploy")
    pgUser := env.Getenv("PGUSER", "gitdeploy")
    pgPassword := env.Getenv("PGPASSWORD", "changeme")
    pgPath := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgDatabase)
    host := env.Getenv("HOST", "")
    port := env.Getenv("PORT", "7654")
    listen := fmt.Sprintf("%s:%s", host, port)

    _ = initPg(pgPath)

    checkDependencies()

    // SSE broker
    sseBroker := sse.New()
    sseBroker.Start()

    // Mux router
    r := mux.NewRouter()
    r.HandleFunc("/deployments", deployWrapperAction(sseBroker)).Methods("POST")
    r.HandleFunc("/deployments", func(w http.ResponseWriter, r *http.Request) {
        // TODO This should be done differently. In some middleware e.g.
        w.Header().Add("Access-Control-Allow-Methods", "POST")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
    }).Methods("OPTIONS")
    r.HandleFunc("/deployments/{app:.+}/events", sseBroker.EventHandler).Methods("GET")
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./app/")))
    http.Handle("/", r)

    go cleanApps()

    log.Printf("Listening on %s", listen)
    log.Fatal(http.ListenAndServe(listen, nil))
}

func initPg(pgPath string) *sql.DB {
    db, err := sql.Open("postgres", pgPath)
    if err != nil {
        log.Fatal(err)
    }
    dat, err := ioutil.ReadFile("./schema.sql")
    if err != nil {
        log.Fatal(err)
    }
    q := string(dat)
    _, err = db.Exec(q)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

func cleanApps() {
    for {
        log.Println("Cleaning up apps...")
        time.Sleep(60 * time.Second)
    }
}

func checkDependencies() {
    go func() {
        if !github.Exists() {
            log.Fatal("Git CLI is required but not installed or not in path.")
        }
    }()
    go func() {
        if !flynn.Exists() {
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

func deployWrapperAction(sseBroker *sse.Broker) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        deployAction(w, r, sseBroker)
    }
}

func deployAction(w http.ResponseWriter, r *http.Request, sseBroker *sse.Broker) {
    // TODO This should be done differently. In some middleware e.g.
    w.Header().Add("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    response := responder.New(apiVersion)
    dr := new(deployRequest)
    decoder := json.NewDecoder(r.Body)
    decoder.Decode(dr)
    app := uuid.NewRandom().String()
    c := sseBroker.OpenChannel(app)

    if len(dr.Repository) < 1 {
        re := response.Error(http.StatusBadRequest, "No repository provided.")
        response.Write(w, re)
        return
    }

    url, err := url.Parse(dr.Repository)
    if err != nil {
        re := response.Error(http.StatusInternalServerError, err.Error())
        response.Write(w, re)
        return
    }

    if url.Host != "github.com" {
        re := response.Error(http.StatusBadRequest, "I only support GitHub.")
        response.Write(w, re)
        return
    }

    re := response.Success(struct {
        App   string `json:"app"`
    }{
        App: app,
    })
    response.Write(w, re)

    go func() {
        defer sseBroker.CloseChannel(app)
        dest, err := cloneRepository(dr.Repository, app, c)
        log.Print("Cloning repository was successful");
        c.Messages <- "Cloning repository was successful"
        if err != nil {
            re := response.Error(http.StatusInternalServerError, err.Error())
            response.Write(w, re)
            return
        }
        err = flynn.Deploy(dest, app, c)
        if err != nil {
            re := response.Error(http.StatusInternalServerError, err.Error())
            response.Write(w, re)
            return
        }
    }()
}

func cloneRepository(repository string, app string, c *sse.Channel) (string, error) {
    dest := fmt.Sprintf("%s/%s", os.TempDir(), app)
    if runtime.GOOS == "windows" {
        dest = fmt.Sprintf("%s\\%s", os.TempDir(), app)
    }

    l := fmt.Sprintf("Creating tempdir %s", dest)
    log.Println(l)
    c.Messages <- l
    return dest, github.Clone(repository, dest, c)
}
