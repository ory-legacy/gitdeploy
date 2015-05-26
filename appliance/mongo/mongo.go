package mongo

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/ory-am/gitdeploy/appliance"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/task/flynn"
)

const (
	eventName = "mongodb.attach"
	procName  = "server"
)

type mongo struct {
	f                *flynn.Flynn
	w                *task.WorkerLog
	h                *task.Helper
	app              string
	workingDirectory string
	eventName        string
	env              map[string]string
}

func New(f *flynn.Flynn, w *task.WorkerLog, h *task.Helper) *mongo {
	return &mongo{f, w, h}
}

type ReleaseContainer struct {
	Manifest string
	URL      string
	*task.Helper
}

func (m *mongo) Attach() (id string, w *task.WorkerLog, env map[string]string, err error) {
	var cw, sw, rw *task.WorkerLog

	// Create app
	id = uuid.NewRandom().String()
	port := 27017
	c := &flynn.CreateApp{m.createHelper(id)}
	if cw, err = c.Run(); err != nil {
		return
	}

	// Write manifest
	manifest, err := appliance.CreateManifest(id, "mongod", procName, port, false)
	if err != nil {
		return
	}

	// Release container
	r := flynn.CreateReleaseContainer(manifest, "url://tbd", id, eventName, m.workingDirectory)
	if rw, err = r.Run(); err != nil {
		return
	}

	// Scale app
	s := &flynn.ScaleApp{ProcName: procName, Helper: m.createHelper(id)}
	if sw, err = s.Run(); err != nil {
		return
	}

	db := uuid.NewRandom().String()
	w = append(cw, sw...)
	w = append(w, rw...)
	env = map[string]string{
		m.env["host"]: id + ".discoverd",
		m.env["port"]: fmt.Sprintf("%d", port),
		m.env["db"]:   db,
		m.env["url"]:  "mongodb://" + id + ":" + fmt.Sprintf("%d", port) + "/" + db,
	}
	return
}

func (m *mongo) createHelper(id string) *task.Helper {
	return &task.Helper{
		App:              id,
		EventName:        eventName,
		WorkingDirectory: m.h.WorkingDirectory,
	}
}
