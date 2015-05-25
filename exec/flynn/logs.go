package flynn

import "os/exec"

func (f *Helper) GetLogs() (string, error) {
    if o, err := exec.Command("flynn", "-a", f.App, "log").CombinedOutput(); err != nil {
        return "", err
    } else {
        return string(o), nil
    }
}
