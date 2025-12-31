package dot

import (
	"fmt"
	"runtime/debug"
)

func ShowVCSInfo() {
	info, _ := debug.ReadBuildInfo()
	var revision string
	var time string
	var modified bool

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.time":
			time = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}

	fmt.Println()
	fmt.Printf("Git Commit: %s\n", revision)
	fmt.Printf("Build Time: %s by go version %s\n", time, info.GoVersion)
	fmt.Printf("Dirty Build: %v\n", modified)
	fmt.Println()
}
