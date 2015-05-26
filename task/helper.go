package task

import (
	"bufio"
	"io"
	"os/exec"
)

type Helper struct {
	App              string
	WorkingDirectory string
	EventName        string
}

// ScanPipe scans the pipe and puts results into w
func (h *Helper) ScanPipe(p io.ReadCloser, w *WorkerLog) {
	s := bufio.NewScanner(p)
	for s.Scan() {
		w.Add(h.EventName, s.Text())
	}
}

func (h *Helper) Exec(w *WorkerLog, cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	e.Dir = h.WorkingDirectory

	if stdout, err := e.StdoutPipe(); err != nil {
		w.AddError(h.EventName, err)
		return err
	} else {
		go h.ScanPipe(stdout, w)
	}

	if stderr, err := e.StderrPipe(); err != nil {
		w.AddError(h.EventName, err)
		return err
	} else {
		go h.ScanPipe(stderr, w)
	}

	err := e.Run()
	if err != nil {
		w.AddError(h.EventName, err)
	}
	return err
}
