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
					go func(appliance *storage.Appliance) {
						e := exec.Command("flynn", "-a", appliance.ID, "delete", "-y")
						if out, err := e.CombinedOutput(); err != nil {
							log.Printf("An error occured while cleanup %s. Reason: %s", err.Error(), out)
							return
						}
					}(&appliance)
				}

				go func(app *storage.App) {
					e := exec.Command("flynn", "-a", app.ID, "delete", "-y")
					out, err := e.CombinedOutput()
					reason := strings.Trim(string(out), " \n\r")
					if err != nil && reason == "controller: resource not found" {
						log.Printf("App %s is not known to controller: %s", app.ID, reason)
					} else if err != nil && reason != "controller: resource not found" {
						log.Printf("An error occured while cleanup %s. Reason: %s", err.Error(), out)
						return
					}
					if err := store.KillApp(app); err != nil {
						log.Printf("Error while cleaning up %s: %s", app, err.Error())
					}
				}(app)
			}
		}
		time.Sleep(15 * time.Second)
	}
}
