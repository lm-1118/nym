package app

import (
	"fmt"
	"nym/core/node"
	"strings"

	"github.com/b9o2/tabby"
)

type InstallCommand struct {
	*tabby.BaseApplication
	ma *MainApp
}

func (c *InstallCommand) Detail() (string, string) {
	return "install", "Install the application"
}

func NewInstallCommand(subApps ...tabby.Application) *InstallCommand {
	return &InstallCommand{
		tabby.NewBaseApplication(0, 0, subApps),
		nil,
	}
}

func (c *InstallCommand) Init(parent tabby.Application) error {
	c.SetParam("version", "要安装的 Node.js 版本", tabby.String(""))
	return nil
}

func (i *InstallCommand) Main(args tabby.Arguments) (*tabby.TabbyContainer, error) {
	if args.IsEmpty() {
		fmt.Println("请指定要安装的Node.js版本。例如：nym install -version 18.17.1")
		return nil, nil
	}

	version, ok := args.Get("version").(string)
	if !ok {
		fmt.Println("参数类型错误")
		return nil, nil
	}

	// 移除前缀v（如果有）
	if strings.HasPrefix(version, "v") {
		version = strings.TrimPrefix(version, "v")
	}

	fmt.Printf("正在安装Node.js v%s...\n", version)

	// 1. 检查版本是否有效
	availableVersions, err := node.ListAvailableVersions()
	if err != nil {
		fmt.Println("获取可用版本失败:", err)
		return nil, nil
	}

	versionExists := false
	for _, v := range availableVersions {
		if v == version {
			versionExists = true
			break
		}
	}

	if !versionExists {
		fmt.Printf("版本v%s不存在。可用的最新5个版本:\n", version)
		for i := 0; i < 5 && i < len(availableVersions); i++ {
			fmt.Println("  " + availableVersions[i])
		}
		return nil, nil
	}

	// 2. 创建进度通道和显示进度
	progressChan := make(chan int)
	go func() {
		var lastProgress int
		for progress := range progressChan {
			if progress > lastProgress {
				fmt.Printf("\r下载进度: %d%%", progress)
				lastProgress = progress
			}
		}
		fmt.Println()
	}()

	// 3. 下载版本
	fmt.Println("开始下载...")
	archivePath, err := node.DownloadVersion(version, progressChan)
	close(progressChan)

	if err != nil {
		fmt.Println("\n下载失败:", err)
		return nil, nil
	}

	// 4. 安装版本
	fmt.Println("\n下载完成，正在安装...")
	if err := node.InstallVersion(version, archivePath); err != nil {
		fmt.Println("安装失败:", err)
		return nil, nil
	}

	fmt.Printf("Node.js v%s 安装成功!\n", version)
	fmt.Println("可以使用 'nym use -version " + version + "' 切换到此版本")

	return nil, nil
}
