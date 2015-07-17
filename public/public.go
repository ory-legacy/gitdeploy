package public

import (
    "log"
    "net/http"
    "os"
    "regexp"
)

func HTML5ModeHandler(dir, index string) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        path := dir + r.URL.Path
        if r.URL.Path == "/" {
            path = dir + "/" + index
        }

        if f, err := os.Stat(path); err == nil {
            if !f.IsDir() {
                http.ServeFile(w, r, path)
                return
            } else {
                http.NotFound(w, r)
                return
            }
        }

        pattern := `!\.html|\.js|\.svg|\.css|\.png|\.jpg$`
        if matched, err := regexp.MatchString(pattern, path); err != nil {
            log.Printf("Could not exec regex: %s", err.Error())
        } else if !matched {
            http.ServeFile(w, r, dir + "/" + index)
            return
        } else {
            http.NotFound(w, r)
        }
    }
}
