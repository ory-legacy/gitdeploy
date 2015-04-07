package config

import (
    "gopkg.in/yaml.v2"
    "github.com/ory-am/gitdeploy/sse"
    "fmt"
    "runtime"
    "os"
    "io/ioutil"
)

type Config struct {
    Process map[string]string
    Go Go
}

type Go struct {
    Package string
}

func Parse(dir string, ch *sse.Channel) error {
    ch.Messages <- "Parsing .gitdeploy.yml"

    c := new(Config)
    filename := fmt.Sprintf("%s/%s", dir, ".gitdeploy.yml")
    if runtime.GOOS == "windows" {
        filename = fmt.Sprintf("%s\\%s", dir, ".gitdeploy.yml")
    }

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        fmt.Printf("no such file or directory: %s", filename)
        ch.Messages <- ".gitdeploy.yml not found, skipping"
        return nil
    }

    data, err := ioutil.ReadFile(filename)
    if err != nil {
        ch.Messages <- err.Error()
        return err
    }

    err = yaml.Unmarshal(data, c)
    if err != nil {
        ch.Messages <- err.Error()
        return err
    }

    var m map[string]string
    d, err := yaml.Marshal(&m)
    if err != nil {
        ch.Messages <- err.Error()
        return err
    }
    ch.Messages <- string(d)

    return nil
}