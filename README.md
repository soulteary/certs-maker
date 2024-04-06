# Certs Maker

[![CodeQL](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml) [![Release](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml) [![Docker Image](https://img.shields.io/docker/pulls/soulteary/certs-maker.svg)](https://hub.docker.com/r/soulteary/certs-maker) [![codecov](https://codecov.io/gh/soulteary/certs-maker/branch/main/graph/badge.svg?token=K12L34CSA4)](https://codecov.io/gh/soulteary/certs-maker)

<p style="text-align: center;">
  <a href="README.md" target="_blank">ENGLISH</a> | <a href="README_CN.md">中文文档</a>
</p>

<img src="logo.png">

**Tiny self-signed tool, file size between 1.5MB and 4MB.**

Generate a self-hosted / dev certificate through configuration.

<img src=".github/assets.jpg">
<img src=".github/dockerhub.jpg">

## Quick Start

Generate self-signed certificate supporting `*.lab.com` and `*.data.lab.com`, just "One Click":

```bash
docker run --rm -it -v `pwd`/ssl:/ssl soulteary/certs-maker:v3.4.1 "--CERT_DNS=lab.com,*.lab.com,*.data.lab.com"
# OR use environment:
# docker run --rm -it -v `pwd`/ssl:/ssl -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" soulteary/certs-maker:v3.4.1
```

Check in the `ssl` directory of the execution command directory:

```bash
ssl
├── lab.com.conf
├── lab.com.crt
└── lab.com.key
```

If you prefer to use file configuration, you can use `docker-compose.yml` like this:

```yaml
version: '2'
services:

certs-maker:
    image: soulteary/certs-maker:v3.4.1
    environment:
      - CERT_DNS=lab.com,*.lab.com,*.data.lab.com
    volumes:
      - ./ssl:/ssl
```

Then execute the following command:

```bash
docker-compose up
# OR
# docker compose up
```

If you want the certificate to be more friendly to K8s, you can add the `FOR_K8S` parameter:

```bash
docker run --rm -it -v `pwd`/ssl:/ssl soulteary/certs-maker:v3.4.1 "--CERT_DNS=lab.com,*.lab.com,*.data.lab.com --FOR_K8S=ON"
# OR
# docker run --rm -it -v `pwd`/ssl:/ssl -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" -e "FOR_K8S=ON" soulteary/certs-maker:v3.4.1
```

And K8S friendly compose file:

```yaml
version: '2'
services:

certs-maker:
    image: soulteary/certs-maker:v3.4.1
    environment:
      - CERT_DNS=lab.com,*.lab.com,*.data.lab.com
      - FOR_K8S=ON
    volumes:
      - ./ssl:/ssl
```

If you want the certificate to be more friendly to Firefox, you can add the `FOR_FIREFOX` parameter:

```bash
docker run --rm -it -v `pwd`/ssl:/ssl soulteary/certs-maker:v3.4.1 "--CERT_DNS=lab.com,*.lab.com,*.data.lab.com --FOR_FIREFOX=ON"
# OR
# docker run --rm -it -v `pwd`/ssl:/ssl -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" -e "FOR_FIREFOX=ON" soulteary/certs-maker:v3.4.1
```

And Firefox friendly compose file:

```yaml
version: '2'
services:

certs-maker:
    image: soulteary/certs-maker:v3.4.1
    environment:
      - CERT_DNS=lab.com,*.lab.com,*.data.lab.com
      - FOR_FIREFOX=ON
    volumes:
      - ./ssl:/ssl
```

If you want to further define the information content of the certificate, including the issuing country, province, street, organization name, etc., you can refer to the following document to manually add parameters.

## SSL certificate parameters

You can customize the generated certificate by declaring the environment variables or cli args of docker.

Use in environment variables:

| Parameter | Name | Use in environment variables |
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

Use in Program CLI arguments:

| Parameter | Name | Use in CLI arguments |
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

## Docker Image

[soulteary/certs-maker](https://hub.docker.com/r/soulteary/certs-maker)
