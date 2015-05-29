package flynn

import (
	"github.com/ory-am/gitdeploy/task"
)

func ScaleApp(w task.WorkerLog, app string, pn string) func(w task.WorkerLog) (error) {
	return func(w task.WorkerLog) (error) {
		w.Add("Releasing container...")
		if err := task.Exec(w, "", "flynn", "-a", app, "scale", pn+"=1"); err != nil {
			return err
		}
		return nil
	}
}

func ReleaseContainer(w task.WorkerLog, app, manifest, url string) func(w task.WorkerLog) (error) {
	return func(w task.WorkerLog) (error) {
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
		return w, err
	}
	return w, nil
}

func CreateApp(w task.WorkerLog, app, wd string) func(w task.WorkerLog) (error) {
	return func(w task.WorkerLog) (error) {
		w.Add("Creating app...")
		if err := task.Exec(w, wd, "flynn", "create", "-y", app); err != nil {
			return err
		}
		return nil
	}
}

func ReleaseApp(w task.WorkerLog, wd string) func(w task.WorkerLog) (error) {
	return func(w task.WorkerLog) (error) {
		w.Add("Releasing app...")
		if err := task.Exec(w, wd, "git", "push", "flynn", "master", "--progress"); err != nil {
			return err
		}
		return nil
	}
}
