package exec

import (
    "github.com/ory-am/event"
    "bufio"
    "io"
    "log"
    "fmt"
    gde "github.com/ory-am/gitdeploy/event"
)

// ScanPipe scans the pipe and puts results into w
func ScanPipe(eventName string, p io.ReadCloser, w *WorkerLog) {
    s := bufio.NewScanner(p)
    for s.Scan() {
        w.Add(eventName, s.Text())
    }
}