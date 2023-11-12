package version

import (
	"runtime"
	"runtime/debug"
)

//nolint:gochecknoglobals
var (
	version   string
	revision  string
	branch    string
	timestamp string
)

type BuildVersion struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	Timestamp string `json:"timestamp"`
	GoVersion string `json:"goVersion"`
	GOARCH    string `json:"goarch"`
	GOOS      string `json:"goos"`
}

func ReadBuildVersion() BuildVersion {
	return BuildVersion{
		Version:   Version(),
		Revision:  Revision(),
		Branch:    Branch(),
		Timestamp: Timestamp(),
		GoVersion: runtime.Version(),
		GOARCH:    runtime.GOARCH,
		GOOS:      runtime.GOOS,
	}
}

func Version() string {
	if version != "" {
		return version
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}

	return ""
}

func Revision() string {
	if revision != "" {
		return revision
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return ""
}

func Branch() string {
	return branch
}

func Timestamp() string {
	if timestamp != "" {
		return timestamp
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}

	return ""
}
