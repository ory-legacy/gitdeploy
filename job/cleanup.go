package job

import (
    "github.com/ory-am/gitdeploy/storage"
    "log"
    "os"
    "os/exec"
    "time"
    "fmt"
)

func KillAppsOnHitList(store storage.Storage) {
    for {
        apps, err := store.GetAppKillList()

        if err != nil {
            log.Printf("An error occured while cleanup: %s", err.Error())
        } else {
            for _, app := range apps {
                go func() {
                    fmt.Println([]string{"flynn", "-a", app.ID, "delete", "-y"})
                    e := exec.Command("flynn", "-a", app.ID, "delete", "-y")
                    e.Stderr = os.Stderr
                    e.Stdout = os.Stdout
                    if err = e.Run(); err != nil {
                        log.Printf("An error occured while cleanup: %s", err.Error())
                    } else {
                        store.KillApp(app)
                    }
                }()
            }
        }

        time.Sleep(15 * time.Second)
    }
}
