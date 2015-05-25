package flynn

import (
    "os/exec"
    gexec "github.com/ory-am/gitdeploy/exec"
)

type Helper struct {
    App string
    WorkingDirectory string
    EventName string
}

func (h *Helper) exec(w *gexec.WorkerLog, cmd string, args ...string) error {
    e := exec.Command(cmd, args...)
    e.Dir = h.WorkingDirectory

    if stdout, err := e.StdoutPipe(); err != nil {
        w.AddError(h.EventName, err.Error())
        return err
    } else {
        go gexec.ScanPipe(h.EventName, stdout, w)
    }

    if stderr, err := e.StderrPipe(); err != nil {
        w.AddError(h.EventName, err.Error())
        return err
    } else {
        go gexec.ScanPipe(h.EventName, stderr, w)
    }

    err := e.Run()
    if err != nil {
        w.AddError(h.EventName, err.Error())
    }
    return err
}