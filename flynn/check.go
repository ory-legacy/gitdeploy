package flynn

import "os/exec"

func Exists() bool {
	_, err := exec.LookPath("flynn")
	return err == nil
}
