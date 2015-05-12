package job

import (
    "fmt"
    "github.com/ory-am/gitdeploy/storage"
    "log"
    "os/exec"
    "strings"
    "time"
    "github.com/ory-am/event"
    gde "github.com/ory-am/gitdeploy/event"
)

func KillAppsOnHitList(store storage.Storage) {
    for {
        apps, err := store.FindAppsOnKillList()
        if err != nil {
            log.Printf("Could not fetch kill-list: %s", err)
        } else {
            for _, app := range apps {
                go func() {
                    fmt.Println([]string{"flynn", "-a", app.ID, "delete", "-y"})
                    e := exec.Command("flynn", "-a", app.ID, "delete", "-y")
                    out, err := e.CombinedOutput()
                    reason := strings.Trim(string(out), " \n\r")
                    if reason == "controller: resource not found" {
                        log.Printf("App %s is not known to controller.", app.ID)
                    } else if err != nil {
                        log.Printf("An error occured while cleanup %s. Reason: %s", err.Error(), out)
                        return
                    }
                    store.KillApp(app)
                }()
            }
        }
        time.Sleep(15 * time.Second)
    }
}

func Cleanup(em *event.EventManager, app, sourcePath string) (error) {
    em.Trigger(eventName, gde.New(app, "Cleaning up..."))
    if err := run(em, eventName, app, sourcePath + "/../", "rm", "-rf", sourcePath); err != nil {
        em.Trigger(eventName, gde.New(app, fmt.Sprintf("Could not remove temp file: %s", err.Error())))
        return err
    }
    return nil
}
