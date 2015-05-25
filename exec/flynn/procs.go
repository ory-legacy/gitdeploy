package flynn

import (
    "os/exec"
)

func (f *Helper) GetProcs(app string) (string, error) {
    if o, err := exec.Command("flynn", "-a", f.App, "ps").CombinedOutput(); err != nil {
        return "", err
    } else {
        return string(o), nil
    }
}
