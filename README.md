# Certs Maker

Generate a self-hosted /dev certificate through configuration.

## Usage

Use `Docker`:

```bash
docker run --rm -it -e parameter=... -v `pwd`/certs:/ssl soulteary/certs-maker
```

Or use `docker-compose`:

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
