package flynn

import (
	"github.com/ory-am/gitdeploy/task"
	"os"
)

type KeyAdd struct{ *task.Helper }
type CreateApp struct{ *task.Helper }
type ReleaseApp struct{ *task.Helper }
type ScaleApp struct {
	ProcName string
	*task.Helper
}
type ReleaseContainer struct {
	Manifest string
	URL      string
	*task.Helper
}

func (d *ScaleApp) Run() (task.WorkerLog, error) {
	w := new(task.WorkerLog)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	w.Add(d.EventName, "Releasing container...")
	if err := d.Exec(w, "flynn", "-a", d.App, "scale", d.ProcName+"=1"); err != nil {
		return err
	}
	return nil
}

func (d *ReleaseContainer) Run() (task.WorkerLog, error) {
	w := new(task.WorkerLog)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	w.Add(d.EventName, "Releasing container...")
	if err := d.Exec(w, "flynn", "-a", d.App, "release", "add", "-f", d.Manifest, d.URL); err != nil {
		return err
	}
	return nil
}

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

func CreateReleaseContainer(manifest, url, id, eventName, workingDirectory string) *ReleaseContainer {
	return &ReleaseContainer{
		Manifest: manifest,
		URL:      "tbd",
		&task.Helper{
			App:              id,
			EventName:        eventName,
			WorkingDirectory: workingDirectory,
		},
	}
}
