package dot

import (
	"fmt"
	"runtime/debug"
	"time"
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
	var buildTime string
	var modified bool

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.time":
			// 解析时间并转换为东八区
			if t, err := time.Parse(time.RFC3339, setting.Value); err == nil {
				// 创建东八区时区
				tz := time.FixedZone("CST", 8*3600)
				// 转换为东八区时间并格式化为 RFC3339
				buildTime = t.In(tz).Format(time.RFC3339)
			} else {
				// 解析失败时使用原始时间
				buildTime = setting.Value
			}
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}

	return VCSInfo{
		Revision:  revision,
		Time:      buildTime,
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
