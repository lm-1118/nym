package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/b9o2/tabby"
)

type InitCommand struct {
	*tabby.BaseApplication
	// ma *MainApp 这行代码的作用是：在 InitCommand 结构体中嵌入一个指向 MainApp 的指针字段。
	// 这样做的好处是，如果 InitCommand 需要访问主应用（MainApp）的一些方法或属性时，可以通过 ma 字段来实现。
	// 例如：可以在 InitCommand 的方法中通过 i.ma 调用 MainApp 的功能，实现子命令与主命令之间的数据或功能交互。
	// 但如果当前没有用到，可以暂时不加，等需要用到主应用上下文时再加也可以。
	ma *MainApp
}

func NewInitCommand() *InitCommand {
	return &InitCommand{
		tabby.NewBaseApplication(0, 0, nil),
		nil,
	}
}

func (i *InitCommand) Detail() (string, string) {
	return "init", "初始化nym环境"
}

func (i *InitCommand) Init(parent tabby.Application) error {
	return nil
}

func (i *InitCommand) Main(args tabby.Arguments) (*tabby.TabbyContainer, error) {
	fmt.Println("正在初始化nym环境")

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("获取用户主目录失败", err)
		return nil, nil
	}
	// 在用户主目录下，拼接出 nym 相关的三个路径，分别用于存放 nym 的主目录、版本目录和当前版本的符号链接。
	nymDir := filepath.Join(home, ".nym")
	versionDir := filepath.Join(nymDir, "version")
	currentDir := filepath.Join(nymDir, "current")

	for _, dir := range []string{nymDir, versionDir, currentDir} {
		// 如果目录不存在，则创建目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("创建目录失败", err)
			return nil, nil
		}
	}
	fmt.Println("目录结构已构建")

	// 设置环境变量
	binPath := filepath.Join(currentDir, "bin")

	//根据不同操作系统设置环境变量
	switch runtime.GOOS {
	case "windows":
		if err := addToWindowsPath(binPath); err != nil {
			fmt.Println("添加环境变量失败", err)
			return nil, nil
		}
	case "linux", "darwin":
		if err := addToLinuxPath(binPath); err != nil {
			fmt.Println("添加环境变量失败", err)
			return nil, nil
		}
	default:
		fmt.Println("不支持的操作系统")
		return nil, nil
	}

	fmt.Println("环境变量已设置")
	fmt.Println("\n初始化完成，请重启终端")

	return nil, nil
}

// 添加环境变量到 Windows 系统,对应的命令是 set PATH=%PATH%;%binPath%
func addToWindowsPath(binPath string) error {
	// 获取当前用户的环境变量
	cmd := exec.Command("powershell", "-Command", "[Environment]::GetEnvironmentVariable('PATH','User')")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	currentPath := strings.TrimSpace(string(output))

	if !strings.Contains(currentPath, binPath) {
		fmt.Println("环境变量已存在，跳过添加")
		return nil
	}

	newPath := currentPath
	if currentPath != "" {
		newPath += ";"
	}
	newPath += binPath
	// 设置环境变量
	cmd = exec.Command("powershell", "-Command",
		fmt.Sprintf("[Environment]::SetEnvironmentVariable('PATH','%s','User')", newPath))
	return cmd.Run()

}

// 添加环境变量到 Linux 系统,对应的命令是 export PATH=$PATH:$binPath
func addToLinuxPath(binPath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// 确定配置文件
	var configFile string
	switch runtime.GOOS {
	case "darwin":
		// 检查用户使用的是bash还是zsh
		if _, err := os.Stat(filepath.Join(home, ".zshrc")); err == nil {
			configFile = filepath.Join(home, ".zshrc")
		} else {
			configFile = filepath.Join(home, ".bash_profile")
		}
	case "linux":
		configFile = filepath.Join(home, ".bashrc")
	}

	// 读取配置文件
	content, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// 要添加的行
	exportLine := fmt.Sprintf("\n# nym Node.js 版本管理器\nexport PATH=\"%s:$PATH\"\n", binPath)

	// 检查是否已经配置
	if strings.Contains(string(content), binPath) {
		fmt.Println("PATH 环境变量已包含 nym 路径")
		return nil
	}

	// 追加到配置文件
	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(exportLine)
	return err

}
