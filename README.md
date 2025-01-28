# Nox

Nox是一个功能强大的网络安全扫描工具箱，集成了端口扫描、子域名发现、Web指纹识别、Web爬虫等多种安全测试功能。

## 功能特性

- 端口扫描：支持TCP/UDP端口扫描，可自定义扫描端口范围和并发数
- 子域名发现：基于字典的子域名枚举，支持DNS解析验证
- Web指纹识别：识别Web应用的技术栈、框架、CMS等信息
- Web爬虫：自动爬取网站结构，发现潜在的敏感信息

# Nox
## 安装

```bash
# 克隆项目
git clone https://github.com/seaung/Nox.git

# 进入项目目录
cd Nox

# 使用Makefile构建
# 执行完整构建流程（清理、格式化、测试和构建）
make all

# 仅构建项目
make build

# 运行测试
make test

# 格式化代码
make fmt

# 运行代码静态检查
make vet

# 清理构建文件
make clean

# 跨平台编译
make build-linux    # 编译Linux版本
make build-windows  # 编译Windows版本
make build-darwin   # 编译MacOS版本
make build-all      # 编译所有平台版本

# 安装到GOPATH
make install
```

## 使用方法
### 端口扫描

```bash
# TCP端口扫描
./nox scan -t example.com -p 80,443,8080

# UDP端口扫描
./nox scan -t example.com -p 53,161 --udp
```

## 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目！

```

## 许可证

本项目采用Apache V2.0许可证，详见[LICENSE](LICENSE)文件。
