package git

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func CreateDirectory(app string) (destination string) {
	tempDir := strings.Trim(os.TempDir(), "/\\")
	destination = fmt.Sprintf("/%s/%s", tempDir, app)
	if runtime.GOOS == "windows" {
		destination = fmt.Sprintf("%s\\%s", tempDir, app)
	}
	return
}
