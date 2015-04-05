package flynn

import (
    "log"
    "github.com/ory-am/gitdeploy/sse"
    "os/exec"
    "bufio"
    "os"
)

func Deploy(sourcePath string, app string, c *sse.Channel) error {
    err := os.Chdir(sourcePath)
    if err != nil {
        return err
    }

    log.Println("Starting with deployment...")
    c.Messages <- "Starting with deployment..."

    log.Println("Adding key...")
    c.Messages <- "Adding key..."
    e := exec.Command("flynn", "key", "add")
    err = run(app, c, e)
    if err != nil {
        log.Fatalf("Error: %s", err.Error())
        c.Messages <- err.Error()
        return err
    }

    log.Println("Creating app...")
    c.Messages <- "Creating app..."
    e = exec.Command("flynn", "create", "-y", app)
    err = run(app, c, e)
    if err != nil {
        log.Fatalf("Error: %s", err.Error())
        c.Messages <- err.Error()
        return err
    }

    log.Println("Pushing to git...")
    c.Messages <- "Pushing to git..."
    c.Messages <- "Due to a bug, you won't receive a response until deployment is finished. This might take some time."
    e = exec.Command("git", "push", "flynn", "master", "--progress")
    err = run(app, c, e)
    if err != nil {
        log.Fatalf("Error: %s", err.Error())
        c.Messages <- err.Error()
        return err
    }

    return err
}

func run(app string, c *sse.Channel, e *exec.Cmd) error {
    stdoutPipe(c, e)
    stderrPipe(c, e)
    return e.Run()
}

func stdoutPipe(c *sse.Channel, e *exec.Cmd) {
    p, err := e.StdoutPipe()
    if err != nil {
        log.Fatal(os.Stderr, "Error creating StdoutPipe for cmd", err)
    }
    s := bufio.NewScanner(p)
    go func() {
        for s.Scan() {
            c.Messages <- string(s.Text())
            log.Printf("%s", s.Text())
        }
    }()
}

func stderrPipe(c *sse.Channel, e *exec.Cmd) {
    p, err := e.StderrPipe()
    if err != nil {
        log.Fatal(os.Stderr, "Error creating StderrPipe for cmd", err)
    }
    s := bufio.NewScanner(p)
    go func() {
        for s.Scan() {
            c.Messages <- string(s.Text())
            log.Printf("%s", s.Text())
        }
    }()
}