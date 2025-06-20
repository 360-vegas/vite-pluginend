package detector

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type SystemInfo struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	Version      string `json:"version"`
	FreeSpace    int64  `json:"free_space"`
	HasAdminPerm bool   `json:"has_admin_perm"`
	NetworkOK    bool   `json:"network_ok"`
}

type DependencyStatus struct {
	Name      string `json:"name"`
	Required  string `json:"required"`
	Current   string `json:"current"`
	Installed bool   `json:"installed"`
	Available bool   `json:"available"`
}

type Detector interface {
	DetectSystem() (*SystemInfo, error)
	CheckDependencies() ([]DependencyStatus, error)
	CheckPorts(ports []int) map[int]bool
}

type detector struct{}

func New() Detector {
	return &detector{}
}

func (d *detector) DetectSystem() (*SystemInfo, error) {
	info := &SystemInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// 检测操作系统版本
	switch runtime.GOOS {
	case "windows":
		info.Version = getWindowsVersion()
	case "darwin":
		info.Version = getMacOSVersion()
	case "linux":
		info.Version = getLinuxVersion()
	}

	// 检测可用磁盘空间
	info.FreeSpace = getDiskFreeSpace()

	// 检测管理员权限
	info.HasAdminPerm = hasAdminPermissions()

	// 检测网络连接
	info.NetworkOK = testNetworkConnection()

	return info, nil
}

func (d *detector) CheckDependencies() ([]DependencyStatus, error) {
	dependencies := []DependencyStatus{
		{Name: "Node.js", Required: ">=18.0.0"},
		{Name: "npm", Required: ">=8.0.0"},
		{Name: "Go", Required: ">=1.21.0"},
		{Name: "Git", Required: ">=2.0.0"},
		{Name: "MySQL", Required: ">=8.0.0"},
		{Name: "MongoDB", Required: ">=6.0.0"},
	}

	for i := range dependencies {
		dep := &dependencies[i]
		dep.Current = getVersionCommand(dep.Name)
		dep.Installed = dep.Current != ""
		dep.Available = isPackageAvailable(dep.Name)
	}

	return dependencies, nil
}

func (d *detector) CheckPorts(ports []int) map[int]bool {
	result := make(map[int]bool)
	for _, port := range ports {
		result[port] = isPortAvailable(port)
	}
	return result
}

func getVersionCommand(name string) string {
	var cmd *exec.Cmd

	switch strings.ToLower(name) {
	case "node.js":
		cmd = exec.Command("node", "--version")
	case "npm":
		cmd = exec.Command("npm", "--version")
	case "go":
		cmd = exec.Command("go", "version")
	case "git":
		cmd = exec.Command("git", "--version")
	case "mysql":
		cmd = exec.Command("mysql", "--version")
	case "mongodb":
		cmd = exec.Command("mongod", "--version")
	default:
		return ""
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

func isPackageAvailable(name string) bool {
	switch runtime.GOOS {
	case "windows":
		return isPackageAvailableWindows(name)
	case "darwin":
		return isPackageAvailableMacOS(name)
	case "linux":
		return isPackageAvailableLinux(name)
	default:
		return false
	}
}

func isPackageAvailableWindows(name string) bool {
	// 检查Chocolatey
	if commandExists("choco") {
		cmd := exec.Command("choco", "search", name, "--exact")
		return cmd.Run() == nil
	}
	// 检查Scoop
	if commandExists("scoop") {
		cmd := exec.Command("scoop", "search", name)
		return cmd.Run() == nil
	}
	return true // 可以通过直接下载安装
}

func isPackageAvailableMacOS(name string) bool {
	if commandExists("brew") {
		cmd := exec.Command("brew", "search", name)
		return cmd.Run() == nil
	}
	return true
}

func isPackageAvailableLinux(name string) bool {
	// 检查不同的包管理器
	packageManagers := []string{"apt", "yum", "dnf", "pacman"}
	for _, pm := range packageManagers {
		if commandExists(pm) {
			return true
		}
	}
	return true
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func getWindowsVersion() string {
	cmd := exec.Command("cmd", "/c", "ver")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown Windows"
	}
	return strings.TrimSpace(string(output))
}

func getMacOSVersion() string {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown macOS"
	}
	return "macOS " + strings.TrimSpace(string(output))
}

func getLinuxVersion() string {
	// 尝试读取/etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				return strings.Trim(strings.Split(line, "=")[1], "\"")
			}
		}
	}

	// 尝试其他方法
	cmd := exec.Command("lsb_release", "-d", "-s")
	if output, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(output))
	}

	return "Unknown Linux"
}

func getDiskFreeSpace() int64 {
	// 简化实现，返回固定值表示有足够空间
	// 在实际应用中可以使用golang.org/x/sys包获取精确的磁盘空间
	return 10 * 1024 * 1024 * 1024 // 10GB
}

func hasAdminPermissions() bool {
	switch runtime.GOOS {
	case "windows":
		// 在Windows上检查管理员权限较复杂，这里简化处理
		cmd := exec.Command("net", "session")
		return cmd.Run() == nil
	case "darwin", "linux":
		return os.Geteuid() == 0
	default:
		return false
	}
}

func testNetworkConnection() bool {
	timeout := time.Second * 3
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func isPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}
