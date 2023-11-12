// version package provides build-time version information.
// It is designed to be used via the -ldflags option of go build.
//
// e.g.
//
//	go build -ldflags "-X github.com/kunitsucom/util.go/version.version=${BUILD_VERSION} -X github.com/kunitsucom/util.go/version.revision=${BUILD_REVISION} -X github.com/kunitsucom/util.go/version.branch=${BUILD_BRANCH} -X github.com/kunitsucom/util.go/version.timestamp=${BUILD_TIMESTAMP}"
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
