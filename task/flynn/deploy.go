package flynn

import (
	"github.com/ory-am/gitdeploy/task"
)

func ScaleApp(app string, pn string) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Releasing container...")
		if err := task.Exec(w, "", "flynn", "-a", app, "scale", pn+"=1"); err != nil {
			return err
		}
		return nil
	}
}

func ReleaseContainer(app, manifest, url string) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Releasing container...")
		if err := task.Exec(w, "", "flynn", "-a", app, "release", "add", "-f", manifest, url); err != nil {
			return err
		}
		return nil
	}
}

func AddKey(w task.WorkerLog) error {
	w.Add("Adding key...")
	if err := task.Exec(w, "", "flynn", "key", "add"); err != nil {
		return err
	}
	return nil
}

func CreateApp(app, wd string) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Creating app...")
		if err := task.Exec(w, wd, "flynn", "create", "-y", app); err != nil {
			return err
		}
		return nil
	}
}

func ReleaseApp(wd string) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Releasing app...")
		if err := task.Exec(w, wd, "git", "push", "flynn", "master", "--progress"); err != nil {
			return err
		}
		return nil
	}
}
