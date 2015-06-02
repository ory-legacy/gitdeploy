package task

import (
	"bufio"
	"io"
	"os/exec"
)

// scanPipe scans an .ReaderCloser and puts results into a WorkerLog
func scanPipe(p io.ReadCloser, w WorkerLog) {
	s := bufio.NewScanner(p)
	for s.Scan() {
		w.Add(s.Text())
	}
}

// Exec executes a command and puts the commands output
func Exec(w WorkerLog, wd string, cmd string, args ...string) error {
	e := exec.Command(cmd, args...)
	if len(wd) != 0 {
		e.Dir = wd + "/"
	}

	if stdout, err := e.StdoutPipe(); err != nil {
		w.AddError(err)
		return err
	} else {
		go scanPipe(stdout, w)
	}

	if stderr, err := e.StderrPipe(); err != nil {
		w.AddError(err)
		return err
	} else {
		go scanPipe(stderr, w)
	}

	if err := e.Run(); err != nil {
		w.AddError(err)
		return err
	}
	return nil
}
