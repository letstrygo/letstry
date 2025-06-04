package editors

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

type EditorName string

func (e EditorName) String() string {
	return string(e)
}

type TrackingType string

func (t TrackingType) String() string {
	return string(t)
}

const (
	TrackingTypeFileAccess TrackingType = "file_access"
	TrackingTypeProcess    TrackingType = "process"
)

var AllTrackingTypes = []TrackingType{
	TrackingTypeFileAccess,
	TrackingTypeProcess,
}

func GetTrackingType(value string) (TrackingType, error) {
	for _, t := range AllTrackingTypes {
		if t.String() == value {
			return t, nil
		}
	}

	return "", fmt.Errorf("unknown tracking type: %s", value)
}

type Editor struct {
	Name                EditorName    `json:"name"`
	ExecPath            string        `json:"path"`
	Args                string        `json:"args"`
	ProcessCaptureDelay time.Duration `json:"process_capture_delay"`
	TrackingType        TrackingType  `json:"tracking_type"`
}

func (e Editor) IsInstalled() bool {
	_, err := os.Stat(e.ExecPath)
	return err == nil
}

func (e Editor) GetExecName() string {
	return filepath.Base(e.ExecPath)
}

func (e Editor) FullString() string {
	return fmt.Sprintf("name: %s, location: %s, args: %s", color.BlueString(e.Name.String()), color.YellowString(e.ExecPath), color.GreenString(e.Args))
}

func (e Editor) String() string {
	return color.BlueString(fmt.Sprintf("(%s, %s)", e.Name, e.GetExecName()))
}
