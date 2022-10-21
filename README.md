# Certs Maker

[![CodeQL](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml)

Small self-signed tool, ~ 3MB Size.

Generate a self-hosted /dev certificate through configuration.


<img src="screenshots/docker.png">

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
| Domians | CERT_DNS | `CERT_DNS=lab.com,*.lab.com,*.data.lab.com` | `--CERT_DNS=yourdomain.com` |
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

## Usage

Generate certificate via cli args:

```bash
docker run --rm -it -v `pwd`/certs:/ssl soulteary/certs-maker --FOR_K8S=on
```

OR use `docker-compose`:

```yaml
version: '2'

services:

  certs-maker:
    image: soulteary/certs-maker
    environment:
      - CERT_DNS=a.com;b.com;c.com;*.d.com;
    volumes:
      - ./certs:/ssl
```

OR, Generate certificate via environment variable:

```bash
docker run --rm -it -e parameter=... -v `pwd`/certs:/ssl soulteary/certs-maker
```

OR, Both use cli args and environment variables:

```bash
docker run --rm -it -e parameter=... -v `pwd`/certs:/ssl soulteary/certs-maker --FOR_K8S=on
```

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
