package job

import (
	"github.com/ory-am/gitdeploy/storage"
	"log"
	"os/exec"
	"strings"
	"time"
)

func KillAppsOnHitList(store storage.Storage) {
	for {
		apps, err := store.FindAppsOnKillList()
		if err != nil {
			log.Printf("Could not fetch kill-list: %s", err)
		} else {
			for _, app := range apps {
				for _, appliance := range app.Appliances {
					go func() {
						e := exec.Command("flynn", "-a", appliance.ID, "delete", "-y")
						if out, err := e.CombinedOutput(); err != nil {
							log.Printf("An error occured while cleanup %s. Reason: %s", err.Error(), out)
							return
						}
					}()
				}

				go func() {
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
