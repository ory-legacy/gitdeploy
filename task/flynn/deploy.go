package flynn

import (
	"github.com/ory-am/gitdeploy/task"
	"log"
)

func ScaleApp(app, pn, amount string) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Releasing container...")
		if err := task.Exec(w, "", "flynn", "-a", app, "scale", pn+"="+amount); err != nil {
			return err
		}
		return nil
	}
}

func ReleaseContainer(app, manifest, url string) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Releasing container...")
		log.Println(w, "", "flynn", "-a", app, "release", "add", "-f", manifest, url)
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

func CreateApp(app, wd string, noRemote bool) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) (err error) {
		w.Add("Creating app...")
		if noRemote {
			err = task.Exec(w, wd, "flynn", "create", "-y", `-r ""`,app)
		} else {
			err = task.Exec(w, wd, "flynn", "create", "-y", app)
		}
		if err != nil {
			return err
		}
		return nil
	}
}

func ReleaseApp(wd, app string) func(w task.WorkerLog) error {
	return func(w task.WorkerLog) error {
		w.Add("Releasing app...")
		if err := task.Exec(w, wd, "git", "push", "--progress", "flynn", app + ":master"); err != nil {
			return err
		}
		return nil
	}
}
