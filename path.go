package dot

import (
	"os"
	"os/exec"
	"path/filepath"
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
