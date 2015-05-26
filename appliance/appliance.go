package appliance
import "github.com/ory-am/gitdeploy/task"

type Appliance interface {
	Attach(args... interface{}) (id string, w *task.WorkerLog, env map[string]string, err error)
	Destroy(id string) (err error)
}
