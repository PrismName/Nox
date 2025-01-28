# Nox

Nox是一个功能强大的网络安全扫描工具箱，集成了端口扫描、子域名发现、Web指纹识别、Web爬虫等多种安全测试功能。

## 功能特性

- 端口扫描：支持TCP/UDP端口扫描，可自定义扫描端口范围和并发数
- 子域名发现：基于字典的子域名枚举，支持DNS解析验证
- Web指纹识别：识别Web应用的技术栈、框架、CMS等信息
- Web爬虫：自动爬取网站结构，发现潜在的敏感信息

## 安装

```bash
# 克隆项目
git clone https://github.com/seaung/Nox.git

# 进入项目目录
cd Nox

# 编译安装
go build -o nox cmd/nox.go
```

## 使用方法

### 端口扫描

```bash
# TCP端口扫描
./nox scan -t example.com -p 80,443,8080

# UDP端口扫描
./nox scan -t example.com -p 53,161 --udp
```

### 子域名发现

```bash
# 使用内置字典扫描子域名
./nox subdomain -d example.com

# 指定自定义字典
./nox subdomain -d example.com -w custom_wordlist.txt
```

### Web指纹识别

```bash
# 识别单个目标
./nox finger -u https://example.com
```

### Web爬虫

```bash
# 爬取网站
./nox crawler -u https://example.com
```

## 注意事项

1. 请合法合规使用本工具，仅用于授权的安全测试
2. 建议在使用前先了解相关法律法规
3. 扫描速度过快可能会触发目标的安全防护机制

## 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目！

## 许可证

本项目采用MIT许可证，详见[LICENSE](LICENSE)文件。
