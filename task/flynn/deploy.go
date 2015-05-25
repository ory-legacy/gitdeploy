package flynn

import (
	"github.com/ory-am/gitdeploy/task"
)

type KeyAdd struct{ *task.Helper }
type CreateApp struct{ *task.Helper }
type ReleaseApp struct{ *task.Helper }

func (d *KeyAdd) Run() (task.WorkerLog, error) {
	w := new(task.WorkerLog)
	w.Add(d.EventName, "Adding key...")

	if err := d.Exec(w, "flynn", "key", "add"); err != nil {
		return w, err
	}
	return w, nil
}

func (d *CreateApp) Run() (task.WorkerLog, error) {
	w := new(task.WorkerLog)
	w.Add(d.EventName, "Creating app...")
	if err := d.Exec(w, "flynn", "create", "-y", d.App); err != nil {
		return err
	}
	return nil
}

func (d *ReleaseApp) Run() (task.WorkerLog, error) {
	w := new(task.WorkerLog)
	w.Add(d.EventName, "Releasing app...")
	if err := d.Exec(w, "git", "push", "flynn", "master", "--progress"); err != nil {
		return err
	}
	return nil
}
