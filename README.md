# Nym - Node Version Manager

Nym 是一个用 Go 语言编写的 Node.js 版本管理工具，类似于 nvm（Node Version Manager）。它提供了一个简单而强大的命令行界面来管理不同的 Node.js 版本。

## 项目架构

### 核心组件

1. **命令行框架 (Tabby)**
   - 使用 [tabby](https://github.com/b9o2/tabby) 作为命令行应用框架
   - 提供了统一的命令行界面和交互体验
   - 支持命令的层级结构和自动补全

2. **命令结构**
   - `RootCmd`: 根命令，作为所有子命令的入口点
   - `ListCommand`: 列出可用的 Node.js 版本
   - `UseCommand`: 切换当前使用的 Node.js 版本

### 代码组织

```
nym/
├── app/           # 命令行相关代码
│   ├── root.go    # 根命令定义
│   ├── list.go    # 列表命令实现
│   └── use.go     # 使用命令实现
├── core/      # 内部实现代码
├── util/          # 工具包
└── main.go       # 程序入口
```

## 设计说明

### 为什么使用 Tabby 框架？

1. **模块化设计**
   - 每个命令都是独立的模块，便于维护和扩展
   - 命令之间通过接口解耦，降低代码耦合度

2. **统一的用户体验**
   - 提供一致的命令行界面
   - 支持命令自动补全和帮助信息
   - 统一的错误处理和输出格式

3. **可扩展性**
   - 易于添加新的命令和功能
   - 支持命令的层级结构
   - 便于集成新的功能模块

### 初始化流程

1. 程序从 `main.go` 开始执行
2. 创建根命令 `RootCmd`
3. 使用 Tabby 框架初始化应用
4. 注册子命令（list、use 等）
5. 启动命令行界面

## 使用说明

### 安装

1. **下载 nym 可执行文件**
   ```bash
   # 下载地址(示例)
   # Windows: https://github.com/your-repo/nym/releases/download/v1.0.0/nym.exe
   # macOS/Linux: https://github.com/your-repo/nym/releases/download/v1.0.0/nym
   ```

2. **使文件可执行** (Linux/macOS)
   ```bash
   chmod +x nym
   mv nym /usr/local/bin/   # 或移动到其他在PATH环境变量中的目录
   ```

### 初始化环境

使用 `init` 命令进行初始化，该命令会创建必要的目录结构并将 nym 添加到环境变量中。

```bash
nym init
```

**注意事项**:
- Windows系统下可能需要管理员权限
- 环境变量修改后需要重新打开终端或运行特定命令使其生效
  - Windows: 重新打开命令提示符或PowerShell
  - Linux/macOS: 执行 `source ~/.bashrc` 或 `source ~/.zshrc`

### 列出已安装的Node.js版本

使用 `list` 命令查看所有已安装的Node.js版本，当前使用的版本会用 `=>` 标记。

```bash
nym list
```

**输出示例**:
```
Installed Node.js versions:
=> v20.5.0 (current)
   v18.17.1
   v16.20.2
```

### 切换Node.js版本

使用 `use` 命令切换到已安装的其他Node.js版本。

```bash
nym use -version 18.17.1
```

或者：

```bash
nym use -version v18.17.1  # 带v前缀也可以
```

**注意事项**:
- Windows系统下可能需要管理员权限才能修改符号链接
- 确保指定的版本已经安装
- 切换版本后，新打开的终端会使用新版本的Node.js

### 版本目录结构

nym 使用以下目录结构管理Node.js版本：

```
~/.nym/
  ├── versions/        # 存放所有已安装的Node.js版本
  │     ├── v20.5.0/   # 特定版本目录
  │     ├── v18.17.1/
  │     └── v16.20.2/
  └── current          # 指向当前使用版本的符号链接
```

其中 `~/.nym/current/bin` 目录会被添加到环境变量PATH中，使系统能够找到当前版本的Node.js可执行文件。

### 常见问题排查

1. **权限问题**
   - Windows下如遇到"创建符号链接失败"，请以管理员身份运行命令或开启开发者模式
   - Linux/macOS下可能需要 `sudo` 权限

2. **版本未更新**
   - 确保重新打开终端或刷新环境变量
   - 检查 `nym list` 输出，确认当前版本是否正确切换

## 开发计划

1. 实现基本的版本管理功能
2. 添加版本下载和安装功能
3. 实现版本切换功能
4. 添加配置管理功能
5. 优化用户界面和交互体验

## 快速开始

### 1. 编译 nym

#### Windows
```powershell
build.bat
```

#### Linux/macOS
```bash
go build -o nym main.go
```

### 2. 把 nym 所在目录加入 PATH

#### Windows
```powershell
setx PATH "%PATH%;D:\Murphy\a\nym"
```
或在“系统属性 → 高级 → 环境变量”里手动添加。

#### Linux/macOS
```bash
echo 'export PATH="$PATH:/home/yourname/nym"' >> ~/.bashrc
source ~/.bashrc
```

### 3. 验证

在任意目录下输入：
```bash
nym --help
```
看到帮助信息即表示配置成功！

---

**这样配置后，你可以在任何目录下直接使用 nym 命令，无需每次切换到 nym 所在目录。**
