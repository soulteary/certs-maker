# Certs Maker

[![CodeQL](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml)

**Tiny self-signed tool, ~ 3MB Size.**

Generate a self-hosted / dev certificate through configuration.

<img src="screenshots/docker.png">

## Quick Start

Generate self-signed certificate supporting `*.lab.com` and `*.data.lab.com`, just "One Click":

```bash
docker run --rm -it -v `pwd`/ssl:/ssl soulteary/certs-maker "--CERT_DNS=lab.com,*.lab.com,*.data.lab.com"
# OR use environment:
# docker run --rm -it -v `pwd`/ssl:/ssl -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" soulteary/certs-maker
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
    image: soulteary/certs-maker
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

## SSL certificate parameters

You can customize the generated certificate by declaring the environment variables or cli args of docker.

Use in environment variables:

| Parameter | Name | Use in environment variables |
| ------ | ------ | ------ |
| Country Name | CERT_C | `CERT_C=CN` | `--CERT_C=CN` |
| State Or Province Name | CERT_ST | `CERT_ST=BJ` | `--CERT_ST=BJ` |
| Locality Name | CERT_L | `CERT_L=HD` | `--CERT_L=HD` |
| Organization Name | CERT_O | `CERT_O=Lab` | `--CERT_O=Lab` |
| Organizational Unit Name | CERT_OU | `CERT_OU=Dev` | `--CERT_OU=Dev` |
| Common Name | CERT_CN | `CERT_CN=Hello World` | `--CERT_CN=Hello World` |
| Domians | CERT_DNS | `CERT_DNS=lab.com,*.lab.com,*.data.lab.com` |
| Issue for K8s | FOR_K8S | `FOR_K8S=ON` | `--FOR_K8S=ON` |

Use in Program CLI arguments:

| Parameter | Name | Use in CLI arguments |
| ------ | ------ | ------ |
| Country Name | CERT_C | `--CERT_C=CN` |
| State Or Province Name | CERT_ST | `--CERT_ST=BJ` |
| Locality Name | CERT_L | `--CERT_L=HD` |
| Organization Name | CERT_O | `--CERT_O=Lab` |
| Organizational Unit Name | CERT_OU | `--CERT_OU=Dev` |
| Common Name | CERT_CN | `--CERT_CN=Hello World` |
| Domians | CERT_DNS | `--CERT_DNS=lab.com,*.lab.com,*.data.lab.com` |
| Issue for K8s | FOR_K8S | `--FOR_K8S=ON` |


## Example

Single domain name:

```bash
# docker run --rm -it -e CERT_DNS=domain.com -v `pwd`/certs:/ssl soulteary/certs-maker

User Input: { CERT_DNS: 'domain.com' }
Generating a RSA private key
..............................................................+++++
.......+++++
writing new private key to 'ssl/domain.com.key'
-----

# ls certs
domain.com.conf domain.com.crt  domain.com.key
```

Wildcard domain name:

```bash
Single domain:

```bash
# docker run --rm -it -e CERT_DNS="*.domain.com" -v `pwd`/certs:/ssl soulteary/certs-maker
# or
# docker run --rm -it -e CERT_DNS=\*.domain.com -v `pwd`/certs:/ssl soulteary/certs-maker

User Input: { CERT_DNS: '*.domain.com' }
Generating a RSA private key
..................+++++
.......................................................+++++
writing new private key to 'ssl/*.domain.com.key'
-----

# ls certs
*.domain.com.conf *.domain.com.crt  *.domain.com.key
```

Multiple domain names:

```bash
Single domain:

```bash
# docker run --rm -it -e CERT_DNS="a.com;*.domain.com;a.c.com" -v `pwd`/certs:/ssl soulteary/certs-maker
# or
# docker run --rm -it -e CERT_DNS=a.com\;\*.domain.com\;a.c.com -v `pwd`/certs:/ssl soulteary/certs-maker

User Input: { CERT_DNS: 'a.com;*.domain.com;a.c.com' }
Generating a RSA private key
...+++++
................................................................................................................................................+++++
writing new private key to 'ssl/a.com.key'
-----

# ls certs
a.com.conf a.com.crt  a.com.key
```

## Docker Image

[soulteary/certs-maker](https://hub.docker.com/r/soulteary/certs-maker)
## LICENSE

[MIT](https://github.com/soulteary/certs-maker/blob/master/LICENSE)
