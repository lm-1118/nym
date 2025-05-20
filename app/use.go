package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/b9o2/tabby"
)

type UseCommand struct {
	*tabby.BaseApplication
	ma *MainApp
}

func NewUseCommand(subApps ...tabby.Application) *UseCommand {
	return &UseCommand{
		tabby.NewBaseApplication(0, 0, subApps),
		nil,
	}
}
func (u *UseCommand) Init(parent tabby.Application) error {
	u.SetParam("version", "要切换的 Node.js 版本", tabby.String(""))
	return nil
}

func (u *UseCommand) Detail() (string, string) {
	return "use", "use version nodejs"
}

func (u *UseCommand) Main(args tabby.Arguments) (*tabby.TabbyContainer, error) {
	if args.IsEmpty() {
		fmt.Println("请指定要切换的 Node.js 版本号")
		return nil, nil
	}
	version, ok := args.Get("version").(string)
	if !ok {
		fmt.Println("参数类型错误")
		return nil, nil
	}

	if strings.HasPrefix(version, "v") {
		version = strings.TrimPrefix(version, "v")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("无法获取用户主目录：", err)
		return nil, nil
	}

	versionsDir := filepath.Join(home, ".nym", "versions")
	targetDir := filepath.Join(versionsDir, "v"+version)
	currentLink := filepath.Join(home, ".nym", "current")

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("版本v%s不存在,请安装改版本。\n", version)
		return nil, nil
	}

	if _, err := os.Lstat(currentLink); err == nil {
		if err := os.Remove(currentLink); err != nil {
			fmt.Println("删除旧的链接失败", err)
			return nil, nil
		}
	}

	err = os.Symlink(targetDir, currentLink)
	if err != nil {
		if runtime.GOOS == "windows" {
			fmt.Println("错误：Windows系统创建符号链接需要管理员权限。")
			fmt.Println("解决方法：")
			fmt.Println("1. 以管理员身份运行 PowerShell，再执行命令")
			fmt.Println("2. 或者开启 Windows 的开发者模式（推荐）")
		} else {
			fmt.Println("创建符号链接失败", err)
		}
		return nil, nil
	}

	fmt.Printf("已经切换到Node.js v%s\n", version)

	return nil, nil
}
