package dot

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func GetBinPath() string {
	binFile, _ := exec.LookPath(os.Args[0])
	binPath, _ := filepath.Abs(binFile)

	return binPath
}

func RootPath() (s string) {
	dir, _ := filepath.Abs(filepath.Dir(s))

	return strings.ReplaceAll(dir, "\\", "/")
}

// GetPackageDir 获取当前函数所在包的目录
func GetPackageDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}
