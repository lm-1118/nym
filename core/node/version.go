package node

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	NodeRootDir = "~/.nym"
	VersionsDir = "~/.nym/versions"
	CurrentLink = "~/.nym/current"
)

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}

func ListInstalledVersions() ([]string, error) {
	versionsDir := ExpandPath(VersionsDir)
	files, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "v") {
			versions = append(versions, strings.TrimPrefix(file.Name(), "v"))
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i] > versions[j]
	})

	return versions, nil
}

func GetCurrentVersion() (string, error) {
	currentLink := ExpandPath(CurrentLink)
	target, err := os.Readlink(currentLink) //读取符号链接,获取它指向的目标路径
	if err != nil {
		return "", err
	}

	versionDir := filepath.Base(target)             //从路径中提取最后一部分
	return strings.TrimPrefix(versionDir, "v"), nil //去掉目录名开头的v,只保留纯版本号
}
