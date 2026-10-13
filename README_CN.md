# Certs Maker / 自制证书工具

[![Go Test Coverage](./.github/coverage.svg)](./.github/go-test-report.md) [![Go Report Card](./.github/goreportcard.svg)](./.github/goreportcard-report.md) [![Go Version](https://img.shields.io/github/go-mod/go-version/soulteary/certs-maker)](./go.mod) [![Release](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml) [![Docker Image](https://img.shields.io/docker/pulls/soulteary/certs-maker.svg)](https://hub.docker.com/r/soulteary/certs-maker)

[English](README.md) | [简体中文](README_CN.md)

<img src="logo.png" alt="Certs Maker 标志" width="160">

一个用于开发和自托管服务的证书生成工具。支持 Docker，以及 macOS、Linux 原生二进制；可为多个域名、通配符域名和 IP 地址生成 PEM、DER 文件。

证书需要客户端显式信任。用于私有生产环境时，请按自己的 PKI 策略管理信任、私钥、部署和续期。工具负责生成文件，不会自动安装证书或续期。

## 快速上手

先安装并启动 Docker，然后执行：

```bash
docker run --rm -v "$(pwd)/ssl:/ssl" \
  -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" \
  soulteary/certs-maker:v3.8.2
```

文件会写入宿主机当前目录下的 `./ssl`：

```text
ssl/
├── lab.com.conf
├── lab.com.der.crt
├── lab.com.der.key
├── lab.com.pem.crt
└── lab.com.pem.key
```

接受 PEM 格式的服务器通常使用 `lab.com.pem.crt` 与 `lab.com.pem.key` 这一对文件。DER 文件以二进制格式保存同一证书和私钥。请妥善保管私钥。

用 OpenSSL 检查证书的签发信息、有效期和 SAN：

```bash
openssl x509 -in ssl/lab.com.pem.crt -noout -subject -issuer -dates -ext subjectAltName
```

生成文件并不会让浏览器自动信任证书。[证书使用与信任指南（英文）](docs/certificates.md)说明了验证、客户端信任、各模式的输出文件，以及 Kubernetes 用法。

示例固定使用 v3.8.2 镜像，包含自定义有效期、Firefox SAN 和 CLI 参数优先级修复。旧版 v3.8.1 的限制见[版本兼容说明（英文）](docs/installation.md#release-compatibility)。

## Docker Compose

将以下配置保存为 `docker-compose.yml`，或使用仓库中的同名文件：

```yaml
services:
  certs-maker:
    image: soulteary/certs-maker:v3.8.2
    environment:
      CERT_DNS: "lab.com,*.lab.com,*.data.lab.com"
    volumes:
      - ./ssl:/ssl
```

执行一次，结束后移除容器：

```bash
docker compose run --rm certs-maker
```

旧版 Compose 可使用 `docker-compose run --rm certs-maker`。

## 证书模式

- **默认模式：** 生成带有 `CA:TRUE`、指定 SAN 的自签名证书，以及 PEM、DER 文件。
- **Kubernetes（`FOR_K8S=ON`）：** 使用 `CA:FALSE` 服务端证书扩展，在 SAN 中添加 `*` 和 `localhost`，并在默认文件名前缀中添加 `.k8s`。Kubernetes Secret 需要自行创建和部署。
- **Firefox（`FOR_FIREFOX=ON`）：** 生成本地根 CA 和单独签发的服务端证书。需要信任它的客户端应导入公开的根 CA 证书。

两个开关同时启用时，Kubernetes 模式优先。[模式、输出文件与验证说明（英文）](docs/certificates.md#certificate-modes)解释了应部署哪份证书、信任哪份证书。

## 配置

每个选项都支持环境变量及同名 CLI 参数。每个参数必须独立传入；带空格或通配符的值需要加引号：

```bash
./certs-maker "--CERT_DNS=lab.com,127.0.0.1" "--CERT_CN=Lab Server" --FOR_K8S=ON
```

[配置参数参考（英文）](docs/configuration.md)集中列出默认值、接受的值、优先级、文件命名、输出路径限制，以及文件所有者设置。

## 原生安装与源码构建

[GitHub Releases](https://github.com/soulteary/certs-maker/releases)提供 macOS、Linux 二进制和 SHA-256 校验文件。原生运行需要 `PATH` 中有 `openssl` 和 `sh`；Docker 镜像已包含这两个依赖。

[安装指南（英文）](docs/installation.md)包含平台选择、校验、源码构建和本地镜像构建。文件及镜像体积随版本和架构变化；下载时的压缩体积与解压后二进制、本地镜像体积不同。

## 常见问题

- **浏览器仍提示证书不可信：** 配置客户端信任并验证域名，参见[信任与验证（英文）](docs/certificates.md#client-trust)。
- **找不到输出文件：** 挂载将宿主机 `./ssl` 映射到容器内 `/ssl`。`DIR` 改变的是容器内路径，参见[输出路径（英文）](docs/configuration.md#output-paths-and-file-names)。
- **能重复执行吗：** 可以，但同名前缀的文件会被覆盖，私钥也会重新生成。先移走原文件或使用新前缀。
- **v3.8.1 为什么忽略自定义有效期，或 Firefox 证书验证失败：** 该版本未包含对应修复。请升级至 v3.8.2，或使用[源码构建（英文）](docs/installation.md#build-from-source)。

## 相关资料

- [《只有 3MB 的自签名证书制作 Docker 工具镜像：Certs Maker》](https://soulteary.com/2022/10/22/make-docker-tools-image-with-only-3md-self-signed-certificate-certs-maker.html)
- [《如何制作和使用自签名证书》](https://soulteary.com/2021/02/06/how-to-make-and-use-a-self-signed-certificate.html)
- [Docker Hub](https://hub.docker.com/r/soulteary/certs-maker)
- [安全政策](SECURITY.md) · [MIT 许可证](LICENSE)
