package editors

import (
	"fmt"
	"os/user"
	"path/filepath"
	"runtime"
	"time"
)

const (
	EditorNameVSCode EditorName = "vscode"
)

func VSCodeEditor() Editor {
	currentUser, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("failed to get current user: %v", err))
	}

	var vsCodePath string
	switch os := runtime.GOOS; os {
	case "darwin":
		vsCodePath = filepath.Join("/", "Applications", "Visual Studio Code.app", "Contents", "Resources", "app", "bin", "code")
	case "linux":
		vsCodePath = filepath.Join(currentUser.HomeDir, "bin", "code")
	case "windows":
		vsCodePath = filepath.Join(currentUser.HomeDir, "AppData", "Local", "Programs", "Microsoft VS Code", "Code.exe")
	}

	return Editor{
		Name:                EditorNameVSCode,
		ExecPath:            vsCodePath,
		Args:                "-n",
		ProcessCaptureDelay: time.Second * 5,
		TrackingType:        TrackingTypeFileAccess,
	}
}
