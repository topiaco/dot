package dot

import (
	"fmt"
	"runtime/debug"
)

// VCSInfo 包含版本控制系统信息
type VCSInfo struct {
	Revision  string
	Time      string
	GoVersion string
	Modified  bool
}

// String 返回格式化的版本信息字符串
func (v VCSInfo) String() string {
	return fmt.Sprintf("\nGit Commit: %s\nBuild Time: %s by go version %s\nDirty Build: %v\n",
		v.Revision, v.Time, v.GoVersion, v.Modified)
}

// extractVCSInfo 从 build info 中提取版本控制信息
func extractVCSInfo() VCSInfo {
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

	return VCSInfo{
		Revision:  revision,
		Time:      time,
		GoVersion: info.GoVersion,
		Modified:  modified,
	}
}

// ShowVCSInfo 打印版本信息
func ShowVCSInfo() {
	fmt.Print(extractVCSInfo().String())
}

// GetVCSInfo 返回 VCSInfo 结构体
func GetVCSInfo() VCSInfo {
	return extractVCSInfo()
}
