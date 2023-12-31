# Certs Maker / 自制证书工具

[![CodeQL](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml) [![Release](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml) [![Docker Image](https://img.shields.io/docker/pulls/soulteary/certs-maker.svg)](https://hub.docker.com/r/soulteary/certs-maker) [![codecov](https://codecov.io/gh/soulteary/certs-maker/branch/main/graph/badge.svg?token=K12L34CSA4)](https://codecov.io/gh/soulteary/certs-maker)

<img src="logo.png">

一个小巧的 SSL 证书生成工具（Docker 工具镜像），静态文件尺寸 1.5MB 左右，容器镜像尺寸 4MB 左右。

你可以使用它快速生成需要的自签名证书，用于生产或开发场景。

<img src=".github/assets.jpg">
<img src=".github/dockerhub.jpg">

## 快速上手

如果你本地已经安装好 Docker 或者 CTR，那么可以通过一条命令快速生成包含 `*.lab.com` 和 `*.data.lab.com` 的证书：

```bash
docker run --rm -it -v `pwd`/ssl:/ssl soulteary/certs-maker:v3.3.0 "--CERT_DNS=lab.com,*.lab.com,*.data.lab.com"
# 如果你希望使用 ENV 来调整生成证书的参数
# docker run --rm -it -v `pwd`/ssl:/ssl -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" soulteary/certs-maker:v3.3.0
```

在命令执行完毕之后，我们检查执行命令的 `ssl` 就能看到生成的证书文件啦：

```bash
ssl
├── lab.com.conf
├── lab.com.crt
└── lab.com.key
```

如果你更喜欢使用配置文件来生成证书，可以使用下面这个 `docker-compose.yml`：

```yaml
version: '2'
services:

certs-maker:
    image: soulteary/certs-maker:v3.3.0
    environment:
      - CERT_DNS=lab.com,*.lab.com,*.data.lab.com
    volumes:
      - ./ssl:/ssl
```

接着，执行下面的命令：

```bash
docker-compose up
# 或者在新版的 compose 中使用下面的命令
# docker compose up
```

如果你希望生成证书对 K8s 使用体验更友好，可以添加 `FOR_K8S` 参数：

```bash
docker run --rm -it -v `pwd`/ssl:/ssl soulteary/certs-maker:v3.3.0 "--CERT_DNS=lab.com,*.lab.com,*.data.lab.com --FOR_K8S=ON"
# 或
# docker run --rm -it -v `pwd`/ssl:/ssl -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" -e "FOR_K8S=ON" soulteary/certs-maker:v3.3.0
```

当然，这里也有使用 `FOR_K8S` 参数的 `compose` 配置文件：

```yaml
version: '2'
services:

certs-maker:
    image: soulteary/certs-maker:v3.3.0
    environment:
      - CERT_DNS=lab.com,*.lab.com,*.data.lab.com
      - FOR_K8S=ON
    volumes:
      - ./ssl:/ssl
```

如果你希望生成证书对 Firefox 的使用体验更友好，可以添加 `FOR_FIREFOX` 参数：

```bash
docker run --rm -it -v `pwd`/ssl:/ssl soulteary/certs-maker:v3.3.0 "--CERT_DNS=lab.com,*.lab.com,*.data.lab.com --FOR_FIREFOX=ON"
# 或
# docker run --rm -it -v `pwd`/ssl:/ssl -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" -e "FOR_FIREFOX=ON" soulteary/certs-maker:v3.3.0
```

当然，这里也有使用 `FOR_FIREFOX` 参数的 `compose` 配置文件：

```yaml
version: '2'
services:

certs-maker:
    image: soulteary/certs-maker:v3.3.0
    environment:
      - CERT_DNS=lab.com,*.lab.com,*.data.lab.com
      - FOR_FIREFOX=ON
    volumes:
      - ./ssl:/ssl
```

如果你希望调整生成证书文件的基础信息（描述信息），诸如签发国家、省份、街道、组织等等，可以参考下面支持的配置参数，进行手动调整。

## SSL 生成工具支持的参数

你可以通过调整环境变量或者 CLI 命令行参数来改变生成的证书。

使用环境变量：

| 类型 | 名称 | 如何在环境变量中使用 |
| ------ | ------ | ------ |
| Country Name | CERT_C | `CERT_C=CN` |
| State Or Province Name | CERT_ST | `CERT_ST=BJ` |
| Locality Name | CERT_L | `CERT_L=HD` |
| Organization Name | CERT_O | `CERT_O=Lab` |
| Organizational Unit Name | CERT_OU | `CERT_OU=Dev` |
| Common Name | CERT_CN | `CERT_CN=Hello World` |
| Domains | CERT_DNS | `CERT_DNS=lab.com,*.lab.com,*.data.lab.com` |
| Issue for K8s | FOR_K8S | `FOR_K8S=ON` |
| Issue for Firefox | FOR_FIREFOX | `FOR_FIREFOX=ON` |
| File Owner User | USER | `USER=ubuntu` |
| File Owner UID | UID | `UID=1234` |
| File Owner GID | GID | `GID=2345` |


使用命令行参数：

| 类型 | 名称 | 如何在环境变量中使用 |
| ------ | ------ | ------ |
| Country Name | CERT_C | `--CERT_C=CN` |
| State Or Province Name | CERT_ST | `--CERT_ST=BJ` |
| Locality Name | CERT_L | `--CERT_L=HD` |
| Organization Name | CERT_O | `--CERT_O=Lab` |
| Organizational Unit Name | CERT_OU | `--CERT_OU=Dev` |
| Common Name | CERT_CN | `--CERT_CN=Hello World` |
| Domains | CERT_DNS | `--CERT_DNS=lab.com,*.lab.com,*.data.lab.com` |
| Issue for K8s | FOR_K8S | `--FOR_K8S=ON` |
| Issue for Firefox | FOR_FIREFOX | `--FOR_FIREFOX=ON` |
| File Owner User | USER | `--USER=ubuntu` |
| File Owner UID | UID | `--UID=1234` |
| File Owner GID | GID | `--GID=2345` |

## Docker 镜像发布地址

[soulteary/certs-maker](https://hub.docker.com/r/soulteary/certs-maker)

## 相关文档教程

- [《只有 3MB 的自签名证书制作 Docker 工具镜像：Certs Maker》](https://soulteary.com/2022/10/22/make-docker-tools-image-with-only-3md-self-signed-certificate-certs-maker.html)
- [《如何制作和使用自签名证书》](https://soulteary.com/2021/02/06/how-to-make-and-use-a-self-signed-certificate.html)
