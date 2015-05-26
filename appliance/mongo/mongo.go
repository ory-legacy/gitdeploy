package mongo

import (
	"github.com/ory-am/gitdeploy/task/flynn"
	"github.com/ory-am/gitdeploy/task"
	"github.com/ory-am/gitdeploy/Godeps/_workspace/src/code.google.com/p/go-uuid/uuid"
)

const eventName = "mongodb.attach"

type appliance struct {
	f *flynn.Flynn
	w *task.WorkerLog
	h *task.Helper
	app              string
	workingDirectory string
	eventName        string
	env map[string]string
}

func New(f *flynn.Flynn, w *task.WorkerLog, h *task.Helper) *appliance {
	return &appliance{f, w, h}
}

func (a *appliance) Attach(_... interface{}) (id string, w *task.WorkerLog, env map[string]string, err error) {
	var cw, sw *task.WorkerLog
	id = uuid.NewRandom().String()
	c := &flynn.CreateApp{a.createHelper(id)}
	if cw, err = c.Run(); err != nil {
		return
	}
	// TODO "server" should be read from manifest.json
	s := &flynn.ScaleApp{
		ProcName: "server",
		Helper: a.createHelper(id),
	}
	if sw, err = s.Run(); err != nil {
		return
	}
	w = append(cw, sw...)
	db := uuid.NewRandom().String()
	env = map[string]string {
		a.env["host"]: id + ".discoverd",
		a.env["port"]: "27017",
		a.env["db"]: db,
		a.env["url"]: "mongodb://" + id + ":27017/" + db,
	}
	return
}

func (a *appliance) createHelper(id string) *task.Helper {
	return &task.Helper{
		App: id,
		EventName: eventName,
		WorkingDirectory: a.h.WorkingDirectory,
	}
}