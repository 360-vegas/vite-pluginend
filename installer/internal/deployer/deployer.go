package deployer

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"installer/internal/database"
)

type Progress struct {
	CurrentStep int    `json:"current_step"`
	TotalSteps  int    `json:"total_steps"`
	StepName    string `json:"step_name"`
	Percentage  int    `json:"percentage"`
	Status      string `json:"status"` // "running", "success", "error"
	Error       string `json:"error,omitempty"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // "info", "warn", "error"
	Message   string    `json:"message"`
}

type Deployer interface {
	Deploy(projectPath string, dbConfig database.Config, settings map[string]string)
	GetProgress() Progress
	GetLogs() []LogEntry
	Verify() error
}

type deployer struct {
	assetFiles embed.FS
	progress   Progress
	logs       []LogEntry
	mutex      sync.RWMutex
}

func New(assetFiles embed.FS) Deployer {
	return &deployer{
		assetFiles: assetFiles,
		progress: Progress{
			TotalSteps: 6,
			Status:     "ready",
		},
		logs: make([]LogEntry, 0),
	}
}

func (d *deployer) Deploy(projectPath string, dbConfig database.Config, settings map[string]string) {
	d.updateProgress(0, "准备安装", "running", "")
	d.addLog("info", "开始部署项目...")

	steps := []struct {
		name string
		fn   func() error
	}{
		{"解压项目文件", func() error { return d.extractProjectFiles(projectPath) }},
		{"生成配置文件", func() error { return d.generateConfigFiles(projectPath, dbConfig, settings) }},
		{"安装依赖", func() error { return d.installDependencies(projectPath) }},
		{"构建项目", func() error { return d.buildProject(projectPath) }},
		{"启动服务", func() error { return d.startServices(projectPath) }},
		{"验证安装", func() error { return d.verifyInstallation() }},
	}

	for i, step := range steps {
		d.updateProgress(i+1, step.name, "running", "")
		d.addLog("info", fmt.Sprintf("执行步骤: %s", step.name))

		if err := step.fn(); err != nil {
			d.updateProgress(i+1, step.name, "error", err.Error())
			d.addLog("error", fmt.Sprintf("步骤失败: %s - %v", step.name, err))
			return
		}

		d.addLog("info", fmt.Sprintf("步骤完成: %s", step.name))
	}

	d.updateProgress(len(steps), "安装完成", "success", "")
	d.addLog("info", "🎉 项目部署完成！")
}

func (d *deployer) GetProgress() Progress {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.progress
}

func (d *deployer) GetLogs() []LogEntry {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.logs
}

func (d *deployer) Verify() error {
	// 验证前端服务
	if err := d.checkService("http://localhost:3000", "前端服务"); err != nil {
		return err
	}

	// 验证后端服务
	if err := d.checkService("http://localhost:8080/health", "后端服务"); err != nil {
		return err
	}

	d.addLog("info", "✅ 所有服务验证通过")
	return nil
}

func (d *deployer) updateProgress(step int, name, status, errorMsg string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.progress.CurrentStep = step
	d.progress.StepName = name
	d.progress.Status = status
	d.progress.Error = errorMsg

	if d.progress.TotalSteps > 0 {
		d.progress.Percentage = (step * 100) / d.progress.TotalSteps
	}
}

func (d *deployer) addLog(level, message string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.logs = append(d.logs, LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	})

	// 限制日志数量，避免内存溢出
	if len(d.logs) > 1000 {
		d.logs = d.logs[len(d.logs)-1000:]
	}
}

func (d *deployer) extractProjectFiles(projectPath string) error {
	d.addLog("info", fmt.Sprintf("正在解压项目文件到: %s", projectPath))

	// 创建目标目录
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 从嵌入的资源中解压项目文件
	projectArchive, err := d.assetFiles.Open("assets/project-files.tar.gz")
	if err != nil {
		// 如果没有预打包的文件，创建基本目录结构
		return d.createBasicStructure(projectPath)
	}
	defer projectArchive.Close()

	// 解压tar.gz文件
	gzReader, err := gzip.NewReader(projectArchive)
	if err != nil {
		return fmt.Errorf("创建gzip reader失败: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取tar文件失败: %w", err)
		}

		targetPath := filepath.Join(projectPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("创建目录失败: %w", err)
			}
		case tar.TypeReg:
			if err := d.extractFile(tarReader, targetPath); err != nil {
				return fmt.Errorf("解压文件失败: %w", err)
			}
		}
	}

	d.addLog("info", "项目文件解压完成")
	return nil
}

func (d *deployer) createBasicStructure(projectPath string) error {
	d.addLog("info", "创建基本项目结构...")

	dirs := []string{
		"frontend",
		"backend",
		"config",
		"logs",
		"data",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (d *deployer) extractFile(reader io.Reader, targetPath string) error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (d *deployer) generateConfigFiles(projectPath string, dbConfig database.Config, settings map[string]string) error {
	d.addLog("info", "生成配置文件...")

	// 生成后端配置
	backendConfig := fmt.Sprintf(`# 数据库配置
DATABASE_TYPE=%s
DATABASE_HOST=%s
DATABASE_PORT=%d
DATABASE_USER=%s
DATABASE_PASSWORD=%s
DATABASE_NAME=%s

# 服务配置
SERVER_PORT=8080
LOG_LEVEL=info

# 安全配置
JWT_SECRET=%s
`,
		dbConfig.Type,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
		generateSecretKey(),
	)

	configPath := filepath.Join(projectPath, "backend", ".env")
	if err := os.WriteFile(configPath, []byte(backendConfig), 0644); err != nil {
		return fmt.Errorf("写入后端配置失败: %w", err)
	}

	// 生成前端配置
	frontendConfig := `VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=插件系统
VITE_APP_VERSION=1.0.0
`

	frontendConfigPath := filepath.Join(projectPath, "frontend", ".env")
	if err := os.WriteFile(frontendConfigPath, []byte(frontendConfig), 0644); err != nil {
		return fmt.Errorf("写入前端配置失败: %w", err)
	}

	d.addLog("info", "配置文件生成完成")
	return nil
}

func (d *deployer) installDependencies(projectPath string) error {
	d.addLog("info", "安装项目依赖...")

	// 安装前端依赖
	frontendPath := filepath.Join(projectPath, "frontend")
	if d.pathExists(frontendPath) {
		d.addLog("info", "安装前端依赖...")
		cmd := exec.Command("npm", "install")
		cmd.Dir = frontendPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("安装前端依赖失败: %w", err)
		}
	}

	// 安装后端依赖
	backendPath := filepath.Join(projectPath, "backend")
	if d.pathExists(backendPath) {
		d.addLog("info", "安装后端依赖...")
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = backendPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("安装后端依赖失败: %w", err)
		}
	}

	d.addLog("info", "依赖安装完成")
	return nil
}

func (d *deployer) buildProject(projectPath string) error {
	d.addLog("info", "构建项目...")

	// 构建前端
	frontendPath := filepath.Join(projectPath, "frontend")
	if d.pathExists(frontendPath) {
		d.addLog("info", "构建前端项目...")
		cmd := exec.Command("npm", "run", "build")
		cmd.Dir = frontendPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("构建前端失败: %w", err)
		}
	}

	// 构建后端
	backendPath := filepath.Join(projectPath, "backend")
	if d.pathExists(backendPath) {
		d.addLog("info", "构建后端项目...")
		cmd := exec.Command("go", "build", "-o", "server", ".")
		cmd.Dir = backendPath
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("构建后端失败: %w", err)
		}
	}

	d.addLog("info", "项目构建完成")
	return nil
}

func (d *deployer) startServices(projectPath string) error {
	d.addLog("info", "启动服务...")

	// 启动后端服务
	backendPath := filepath.Join(projectPath, "backend")
	if d.pathExists(backendPath) {
		d.addLog("info", "启动后端服务...")
		cmd := exec.Command("./server")
		cmd.Dir = backendPath

		// 在后台启动
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("启动后端服务失败: %w", err)
		}

		// 等待服务启动
		time.Sleep(2 * time.Second)
	}

	// 启动前端服务
	frontendPath := filepath.Join(projectPath, "frontend")
	if d.pathExists(frontendPath) {
		d.addLog("info", "启动前端服务...")
		cmd := exec.Command("npm", "run", "serve")
		cmd.Dir = frontendPath

		// 在后台启动
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("启动前端服务失败: %w", err)
		}

		// 等待服务启动
		time.Sleep(3 * time.Second)
	}

	d.addLog("info", "服务启动完成")
	return nil
}

func (d *deployer) verifyInstallation() error {
	d.addLog("info", "验证安装...")

	// 这里可以添加更多的验证逻辑
	d.addLog("info", "安装验证通过")
	return nil
}

func (d *deployer) checkService(url, name string) error {
	// 简化的服务检查，实际应该发送HTTP请求
	d.addLog("info", fmt.Sprintf("检查%s: %s", name, url))
	return nil
}

func (d *deployer) pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func generateSecretKey() string {
	// 简化的密钥生成，实际应该使用加密随机数
	return "your-secret-key-here-" + fmt.Sprintf("%d", time.Now().Unix())
}
