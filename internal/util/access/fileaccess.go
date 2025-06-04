package access

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func IsPathUse(path string) bool {
	switch runtime.GOOS {
	case "windows":
		newPath := fmt.Sprintf("%s-%v", path, time.Now().Unix())
		err := os.Rename(path, newPath)
		if err != nil {
			return true
		}

		_ = os.Rename(newPath, path)
		return false
	default:
		cmd := exec.Command("lsof", path)

		if err := cmd.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				ec := exitError.ExitCode()
				if ec == 1 {
					return false
				}
			}
		}

		return true
	}
}
