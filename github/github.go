package github

import (
	"os/exec"
	"github.com/ory-am/gitdeploy/sse"
	"log"
	"os"
	"bufio"
)

func Exists() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func Clone(source string, destination string, c *sse.Channel) (error) {
	log.Println("Starting git clone...")
	c.Messages <- "Starting git clone..."
	e := exec.Command("git", "clone", "--progress", source, destination)

	er1, err := e.StdoutPipe()
	if err != nil {
		log.Fatal(os.Stderr, "Error creating StdoutPipe for Cmd", err)
	}
	s1 := bufio.NewScanner(er1)
	go func() {
		for s1.Scan() {
			c.Messages <- string(s1.Text())
			log.Printf("Git subcommand: %s", s1.Text())
		}
	}()

	er2, err := e.StderrPipe()
	if err != nil {
		log.Fatal(os.Stderr, "Error creating StdoutPipe for Cmd", err)
	}
	s2 := bufio.NewScanner(er2)
	go func() {
		for s2.Scan() {
			c.Messages <- string(s2.Text())
			log.Printf("Git subcommand failed: %s", s2.Text())
		}
	}()

	return e.Run()
}
