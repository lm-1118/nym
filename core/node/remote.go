package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func GetVersionDownloadURL(version string) string {
	// 移除前缀v（如果有）
	if strings.HasPrefix(version, "v") {
		version = strings.TrimPrefix(version, "v")
	}
	arch := runtime.GOARCH
	os := runtime.GOOS
	var nodeArch, nodeOS string
	switch arch {
	case "amd64":
		nodeArch = "x64"
	case "arm64":
		nodeArch = "arm64"
	case "386":
		nodeArch = "x86"
	default:
		nodeArch = arch
	}

	switch os {
	case "windows":
		nodeOS = "win"
		return fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%s-%s-%s.zip", version, version, nodeOS, nodeArch)
	case "darwin":
		nodeOS = "darwin"
		return fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%s-%s-%s.tar.gz", version, version, nodeOS, nodeArch)
	default:
		nodeOS = "linux"
		return fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%s-%s-%s.tar.gz", version, version, nodeOS, nodeArch)
	}

}

func ListAvailableVersions() ([]string, error) {
	//请求Node.js官网获取版本列表
	resp, err := http.Get("https://nodejs.org/dist/index.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//解析JSON响应
	// 这里定义了一个结构体切片，用于存储从Node.js官网获取到的所有版本信息
	var versions []struct {
		Version string `json:"version"` // 对应JSON中的"version"字段
	}
	// 使用json.NewDecoder(resp.Body).Decode(&versions)方法
	// 作用：从HTTP响应体中读取JSON数据，并将其解析（反序列化）到上面定义的versions切片中
	// 如果解析过程中出错（比如格式不对），就返回错误
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, err
	}

	//提取版本号
	var result []string
	// 切片的常用方法有：append、len、cap、copy、切片截取（a[start:end]）、遍历（for ... range）
	// strings包的常用方法有：TrimPrefix、TrimSuffix、Contains、HasPrefix、HasSuffix、Split、Join、Replace、ToLower、ToUpper、Index、LastIndex
	for _, v := range versions {
		// 去掉每个版本号前面的"v"，比如"v18.17.1"变成"18.17.1"
		result = append(result, strings.TrimPrefix(v.Version, "v"))
	}

	return result, nil
}
