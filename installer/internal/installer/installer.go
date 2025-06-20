package installer

import (
	"fmt"
	"os/exec"
	"runtime"
)

type Manager interface {
	Install(pkg string) error
	IsAvailable(pkg string) bool
}

type manager struct {
	packageManager PackageManager
}

type PackageManager interface {
	Install(pkg string) error
	IsInstalled(pkg string) bool
	IsAvailable(pkg string) bool
}

func New() Manager {
	var pm PackageManager

	switch runtime.GOOS {
	case "windows":
		pm = NewWindowsPackageManager()
	case "darwin":
		pm = NewMacOSPackageManager()
	case "linux":
		pm = NewLinuxPackageManager()
	default:
		pm = &DirectDownloadManager{}
	}

	return &manager{packageManager: pm}
}

func (m *manager) Install(pkg string) error {
	return m.packageManager.Install(pkg)
}

func (m *manager) IsAvailable(pkg string) bool {
	return m.packageManager.IsAvailable(pkg)
}

// Windows包管理器
type WindowsPackageManager struct {
	hasChoco bool
	hasScoop bool
}

func NewWindowsPackageManager() *WindowsPackageManager {
	return &WindowsPackageManager{
		hasChoco: commandExists("choco"),
		hasScoop: commandExists("scoop"),
	}
}

func (w *WindowsPackageManager) Install(pkg string) error {
	if w.hasChoco {
		return runCommand("choco", "install", pkg, "-y")
	} else if w.hasScoop {
		return runCommand("scoop", "install", pkg)
	}

	// 回退到直接下载
	return w.directDownload(pkg)
}

func (w *WindowsPackageManager) IsInstalled(pkg string) bool {
	return commandExists(pkg)
}

func (w *WindowsPackageManager) IsAvailable(pkg string) bool {
	return w.hasChoco || w.hasScoop
}

func (w *WindowsPackageManager) directDownload(pkg string) error {
	downloadMap := map[string]string{
		"nodejs": "https://nodejs.org/dist/v20.10.0/node-v20.10.0-x64.msi",
		"go":     "https://golang.org/dl/go1.21.5.windows-amd64.msi",
		"git":    "https://github.com/git-for-windows/git/releases/download/v2.43.0.windows.1/Git-2.43.0-64-bit.exe",
	}

	url, ok := downloadMap[pkg]
	if !ok {
		return fmt.Errorf("不支持直接下载: %s", pkg)
	}

	return downloadAndInstall(url)
}

// macOS包管理器
type MacOSPackageManager struct {
	hasHomebrew bool
}

func NewMacOSPackageManager() *MacOSPackageManager {
	return &MacOSPackageManager{
		hasHomebrew: commandExists("brew"),
	}
}

func (m *MacOSPackageManager) Install(pkg string) error {
	if !m.hasHomebrew {
		// 安装Homebrew
		if err := m.installHomebrew(); err != nil {
			return err
		}
		m.hasHomebrew = true
	}

	// 映射包名
	brewPkg := mapPackageName(pkg, "macos")
	return runCommand("brew", "install", brewPkg)
}

func (m *MacOSPackageManager) IsInstalled(pkg string) bool {
	return commandExists(pkg)
}

func (m *MacOSPackageManager) IsAvailable(pkg string) bool {
	return m.hasHomebrew || true // Homebrew可以安装
}

func (m *MacOSPackageManager) installHomebrew() error {
	cmd := `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
	return exec.Command("sh", "-c", cmd).Run()
}

// Linux包管理器
type LinuxPackageManager struct {
	packageManager string
}

func NewLinuxPackageManager() *LinuxPackageManager {
	pm := detectLinuxPackageManager()
	return &LinuxPackageManager{packageManager: pm}
}

func (l *LinuxPackageManager) Install(pkg string) error {
	linuxPkg := mapPackageName(pkg, "linux")

	switch l.packageManager {
	case "apt":
		if err := runCommand("apt", "update"); err != nil {
			return err
		}
		return runCommand("apt", "install", "-y", linuxPkg)
	case "yum":
		return runCommand("yum", "install", "-y", linuxPkg)
	case "dnf":
		return runCommand("dnf", "install", "-y", linuxPkg)
	case "pacman":
		return runCommand("pacman", "-S", "--noconfirm", linuxPkg)
	default:
		return fmt.Errorf("不支持的包管理器: %s", l.packageManager)
	}
}

func (l *LinuxPackageManager) IsInstalled(pkg string) bool {
	return commandExists(pkg)
}

func (l *LinuxPackageManager) IsAvailable(pkg string) bool {
	return l.packageManager != ""
}

// 直接下载管理器（备用方案）
type DirectDownloadManager struct{}

func (d *DirectDownloadManager) Install(pkg string) error {
	return fmt.Errorf("直接下载安装暂未实现: %s", pkg)
}

func (d *DirectDownloadManager) IsInstalled(pkg string) bool {
	return commandExists(pkg)
}

func (d *DirectDownloadManager) IsAvailable(pkg string) bool {
	return true
}

// 辅助函数
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func detectLinuxPackageManager() string {
	packageManagers := []string{"apt", "yum", "dnf", "pacman"}
	for _, pm := range packageManagers {
		if commandExists(pm) {
			return pm
		}
	}
	return ""
}

func mapPackageName(pkg, platform string) string {
	// 映射通用包名到平台特定包名
	packageMap := map[string]map[string]string{
		"nodejs": {
			"macos": "node",
			"linux": "nodejs",
		},
		"Node.js": {
			"macos": "node",
			"linux": "nodejs",
		},
		"mysql": {
			"macos": "mysql",
			"linux": "mysql-server",
		},
		"MongoDB": {
			"macos": "mongodb-community",
			"linux": "mongodb",
		},
	}

	if platformMap, ok := packageMap[pkg]; ok {
		if mappedName, ok := platformMap[platform]; ok {
			return mappedName
		}
	}

	return pkg
}

func downloadAndInstall(url string) error {
	// 简化实现，实际应该下载并执行安装程序
	return fmt.Errorf("直接下载功能暂未实现: %s", url)
}
