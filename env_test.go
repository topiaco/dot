package dot

import (
	"os"
	"testing"
)

// TestIsInContainer 测试判断容器环境方法
func TestIsInContainer(t *testing.T) {
	_ = IsInContainer()

	// 测试环境变量模拟
	orig := os.Getenv("KUBERNETES_SERVICE_HOST")
	defer os.Setenv("KUBERNETES_SERVICE_HOST", orig)

	os.Setenv("KUBERNETES_SERVICE_HOST", "10.0.0.1")
	if !IsInContainer() {
		t.Errorf("期望当 KUBERNETES_SERVICE_HOST 存在时返回 true")
	}
}
