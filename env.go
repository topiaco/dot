package dot

import (
	"os"
	"strings"
)

func Getenv(key string, def ...string) string {
	val := os.Getenv(key)
	if val == "" {
		if len(def) > 0 {
			return def[0]
		}
	}

	return val
}

// IsInContainer 判断当前运行环境是否在容器中
func IsInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount"); err == nil {
		return true
	}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		content := string(data)
		if strings.Contains(content, "docker") ||
			strings.Contains(content, "kubepods") ||
			strings.Contains(content, "containerd") ||
			strings.Contains(content, "lxc") {
			return true
		}
	}
	return false
}
