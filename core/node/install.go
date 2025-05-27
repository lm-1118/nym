package node

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// DownloadVersion 用于下载指定版本的 Node.js 安装包，并实时反馈下载进度。
// version: 需要下载的 Node.js 版本号（如 "18.17.1"）
// progressChan: 用于传递下载进度的通道（百分比，0~100）
// 返回值：下载完成后的文件路径，以及可能出现的错误
func DownloadVersion(version string, progressChan chan<- int) (string, error) {
	// 1. 获取下载链接
	url := GetVersionDownloadURL(version)
	fmt.Println("下载URL:", url)

	// 2. 创建临时目录用于存放下载的文件
	tempDir, err := os.MkdirTemp("", "nym-download-")
	if err != nil {
		return "", err // 创建临时目录失败
	}

	// 3. 解析文件名和目标文件路径
	fileName := filepath.Base(url)
	filePath := filepath.Join(tempDir, fileName)

	// 4. 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", err // 创建文件失败
	}
	defer file.Close() // 函数结束时关闭文件

	// 5. 发送 HTTP GET 请求下载文件
	resp, err := http.Get(url)
	if err != nil {
		return "", err // 下载失败
	}
	defer resp.Body.Close() // 函数结束时关闭响应体

	// 6. 获取文件总大小（字节数），用于计算进度
	fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// 7. 创建缓冲区，每次读取 1MB
	buffer := make([]byte, 1024*1024)
	downloaded := 0 // 已下载字节数

	// 8. 循环读取并写入文件，同时上报进度
	for {
		n, err := resp.Body.Read(buffer) // 读取数据到缓冲区
		if n > 0 {
			file.Write(buffer[:n]) // 写入文件
			downloaded += n        // 累加已下载字节数
			// 计算进度百分比并发送到通道
			if fileSize > 0 {
				progress := int(float64(downloaded) / float64(fileSize) * 100)
				progressChan <- progress
			}
		}
		if err != nil {
			if err == io.EOF {
				break // 文件读取完毕，退出循环
			}
			return "", err // 读取过程中出错
		}
	}

	// 下载完成后验证文件
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("无法获取文件信息: %v", err)
	}

	// 检查文件大小
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("下载的文件为空")
	}

	// 如果是ZIP文件，简单验证文件头
	if strings.HasSuffix(filePath, ".zip") {
		f, err := os.Open(filePath)
		if err != nil {
			return "", err
		}
		defer f.Close()

		// 读取文件头
		header := make([]byte, 4)
		if _, err := f.Read(header); err != nil {
			return "", err
		}

		// ZIP文件头通常以 PK\03\04 开始
		if header[0] != 0x50 || header[1] != 0x4B {
			return "", fmt.Errorf("不是有效的ZIP文件")
		}
	}

	return filePath, nil // 返回下载文件的路径
}

// InstallVersion 用于安装指定版本的 Node.js，目前未实现
// version: 需要安装的 Node.js 版本号
// archivePath: 下载的 Node.js 安装包路径
// 返回值：错误信息（如果有）
func InstallVersion(version string, archivePath string) error {
	// 1. 解压安装包
	home, _ := os.UserHomeDir()
	versionsDir := filepath.Join(home, ".nym", "versions")
	targetDir := filepath.Join(versionsDir, "v"+version)

	os.MkdirAll(targetDir, 0755)

	// 2. 解压文件
	switch {
	case strings.HasSuffix(archivePath, ".tar.gz"):
		err := untarFile(archivePath, targetDir)
		if err != nil {
			return err
		}
	case strings.HasSuffix(archivePath, ".zip"):
		err := unzipFile(archivePath, targetDir)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("不支持的文件格式: %s", archivePath)
	}

	os.Remove(archivePath)

	return nil
}

// 解压TAR.GZ文件到目标目录
func untarFile(tarPath string, destDir string) error {
	// 打开压缩文件
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建gzip读取器
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// 创建tar读取器
	tarReader := tar.NewReader(gzipReader)

	// 提取tar文件中的第一级目录名（通常是 node-v12.18.3-darwin-x64 这样的格式）
	var baseDir string
	firstHeader, err := tarReader.Next()
	if err == nil && firstHeader.Typeflag == tar.TypeDir {
		baseDir = firstHeader.Name
	}

	// 重新打开文件并重置读取器
	file.Seek(0, 0)
	gzipReader, _ = gzip.NewReader(file)
	tarReader = tar.NewReader(gzipReader)

	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// 解压所有文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 文件结束
		}
		if err != nil {
			return err
		}

		// 跳过第一级目录，直接提取内容到目标目录
		if baseDir != "" && strings.HasPrefix(header.Name, baseDir+"/") {
			// 计算相对路径
			relPath := strings.TrimPrefix(header.Name, baseDir+"/")
			if relPath == "" {
				continue // 跳过目录本身
			}

			// 计算目标路径
			destPath := filepath.Join(destDir, relPath)

			switch header.Typeflag {
			case tar.TypeDir:
				// 创建目录
				if err := os.MkdirAll(destPath, 0755); err != nil {
					return err
				}
			case tar.TypeReg:
				// 创建文件
				destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
				if err != nil {
					return err
				}

				// 复制内容
				if _, err := io.Copy(destFile, tarReader); err != nil {
					destFile.Close()
					return err
				}
				destFile.Close()
			}
		}
	}

	return nil
}

// 解压ZIP文件到目标目录
func unzipFile(zipPath string, destDir string) error {
	// 打开ZIP文件
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// 提取ZIP文件中的第一级目录名（通常是 node-v12.18.3-win-x64 这样的格式）
	var baseDir string
	if len(reader.File) > 0 {
		firstFile := reader.File[0]
		parts := strings.Split(firstFile.Name, "/")
		if len(parts) > 0 {
			baseDir = parts[0]
		}
	}

	// 解压所有文件
	for _, file := range reader.File {
		// 跳过第一级目录，直接提取内容到目标目录
		if baseDir != "" && strings.HasPrefix(file.Name, baseDir+"/") {
			// 计算相对路径
			relPath := strings.TrimPrefix(file.Name, baseDir+"/")
			if relPath == "" {
				continue // 跳过目录本身
			}

			// 计算目标路径
			destPath := filepath.Join(destDir, relPath)

			if file.FileInfo().IsDir() {
				// 创建目录
				if err := os.MkdirAll(destPath, file.Mode()); err != nil {
					return err
				}
			} else {
				// 创建文件
				destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
				if err != nil {
					return err
				}

				// 打开源文件
				srcFile, err := file.Open()
				if err != nil {
					destFile.Close()
					return err
				}

				// 复制内容
				_, err = io.Copy(destFile, srcFile)
				srcFile.Close()
				destFile.Close()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
