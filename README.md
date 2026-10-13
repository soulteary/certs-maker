# Certs Maker

Generate TLS certificates for development and self-hosted services with a single command.

[![Go Test Coverage](./.github/coverage.svg)](./.github/go-test-report.md) [![Go Report Card](./.github/goreportcard.svg)](./.github/goreportcard-report.md) [![Go Version](https://img.shields.io/github/go-mod/go-version/soulteary/certs-maker)](./go.mod) [![Release](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml) [![Docker Image](https://img.shields.io/docker/pulls/soulteary/certs-maker.svg)](https://hub.docker.com/r/soulteary/certs-maker)

[English](README.md) | [简体中文](README_CN.md)

<img src="logo.png" alt="Certs Maker logo" width="160">

Certs Maker turns a list of DNS names and IP addresses into certificate files you can use with a web server, reverse proxy, or local service. It runs as a one-shot Docker container or a native macOS/Linux program and uses OpenSSL to generate RSA 2048-bit keys and SHA-256 signatures.

- Include multiple domains, wildcard domains, IPv4, and IPv6 addresses in one certificate.
- Export a certificate and private key in both PEM and DER formats.
- Choose a self-signed certificate or a local root CA with a signed server certificate.
- Configure names, validity, and output through environment variables or CLI flags.

Certificates require explicit client trust. You manage server configuration, private keys, client trust, and renewal; Certs Maker generates the files.

## Quick start

With Docker installed and running, generate a certificate for `lab.com`, its direct subdomains, and `127.0.0.1`:

```bash
docker run --rm -v "$(pwd)/ssl:/ssl" \
  -e "CERT_DNS=lab.com,*.lab.com,127.0.0.1" \
  soulteary/certs-maker:v3.8.1
```

The container exits after generation. Your files are saved to `./ssl` on the host:

```text
ssl/
├── lab.com.conf
├── lab.com.der.crt
├── lab.com.der.key
├── lab.com.pem.crt
└── lab.com.pem.key
```

For a server that accepts PEM files, configure:

- **Certificate:** `ssl/lab.com.pem.crt`
- **Private key:** `ssl/lab.com.pem.key`

The DER pair contains the same certificate and key in binary form. Private keys are unencrypted; restrict access to them.

**Version note:** the published v3.8.1 image supports this basic example, but predates the custom-validity, Firefox SAN, and CLI-precedence fixes on `main`. Build [current source](#build-from-source) to use those fixes. See [release compatibility](docs/installation.md#release-compatibility) for details.

### Docker Compose

Save the following as `docker-compose.yml`:

```yaml
services:
  certs-maker:
    image: soulteary/certs-maker:v3.8.1
    environment:
      CERT_DNS: "lab.com,*.lab.com,127.0.0.1"
    volumes:
      - ./ssl:/ssl
```

Run it once and remove the container afterward:

```bash
docker compose run --rm certs-maker
```

The bind mount maps host `./ssl` to container `/ssl`. To use another host directory, change the source of the mount, for example `./certificates:/ssl`.

## Build from source

Current source includes the fixes described above. You need the Go version required by [`go.mod`](go.mod), currently **Go 1.27.2**, plus `openssl` and `sh` in `PATH`.

```bash
git clone https://github.com/soulteary/certs-maker.git
cd certs-maker
go build -o certs-maker .
./certs-maker "--CERT_DNS=lab.com,*.lab.com,127.0.0.1" --EXPIRE_DAYS=30
```

This generates a certificate valid for 30 days in `./ssl`.

Prebuilt macOS and Linux archives are available from [GitHub Releases](https://github.com/soulteary/certs-maker/releases). Native binaries also require OpenSSL and a shell. The [installation guide](docs/installation.md) covers platform selection, checksums, and packaging current source into a local Docker image.

## Choose a certificate mode

The following behavior applies to current source. For Firefox mode and custom validity, use a source build until a published release includes the fixes.

| Mode | Setting | Certificate and output |
| --- | --- | --- |
| Default | Both mode flags off | Self-signed certificate with `CA:TRUE`; five PEM/DER/configuration files. |
| Kubernetes | `FOR_K8S=ON` | Self-signed server certificate with `CA:FALSE`; default file prefix gains `.k8s`. |
| Firefox / local CA | `FOR_FIREFOX=ON` | Local root CA plus a separately signed server certificate with `CA:FALSE` and DNS/IP SANs. |

### Kubernetes

```bash
./certs-maker "--CERT_DNS=lab.com,*.lab.com,127.0.0.1" \
  --FOR_K8S=ON --DIR=ssl-k8s
```

Use `ssl-k8s/lab.com.k8s.pem.crt` and `ssl-k8s/lab.com.k8s.pem.key` for your server or TLS Secret. This mode also adds `localhost` and `*` SAN entries. The `*` entry does not make the certificate valid for every hostname: always include the actual DNS names and IP addresses clients use.

You create and configure the Kubernetes Secret or Ingress yourself. See the [Kubernetes usage example](docs/certificates.md#kubernetes-mode).

### Firefox / local CA

```bash
./certs-maker "--CERT_DNS=lab.com,*.lab.com,127.0.0.1" \
  --FOR_FIREFOX=ON --EXPIRE_DAYS=30 --DIR=ssl-firefox
```

Along with the ordinary five files, this creates `.rootCA.pem`, `.rootCA.key`, `.pem.csr`, and `.rootCA.srl` files.

- **Server:** deploy `ssl-firefox/lab.com.pem.crt` and `ssl-firefox/lab.com.pem.key`.
- **Clients:** trust the public `ssl-firefox/lab.com.rootCA.pem` certificate.
- **Private material:** keep both `.pem.key` and `.rootCA.key` private.

Firefox users can import the public root certificate under **Authorities** in the certificate manager and enable trust for website identification. Other clients may use a different system, application, or container trust store. See the [client trust guide](docs/certificates.md#client-trust).

If both mode flags are enabled, Kubernetes mode takes precedence and no Firefox root CA is generated.

## Verify before deployment

These verification examples use OpenSSL 3.x. On macOS, the bundled `openssl` may be LibreSSL and may not support these options; use an OpenSSL 3.x binary for the commands below.

Inspect the certificate’s dates and Subject Alternative Names (SANs):

```bash
openssl x509 -in ssl/lab.com.pem.crt \
  -noout -subject -issuer -dates -ext subjectAltName
```

For the default self-signed certificate, verify a hostname against that certificate:

```bash
openssl verify -CAfile ssl/lab.com.pem.crt \
  -verify_hostname lab.com ssl/lab.com.pem.crt
```

For Firefox mode, verify the server certificate against its root CA:

```bash
openssl verify -CAfile ssl-firefox/lab.com.rootCA.pem \
  -verify_hostname lab.com ssl-firefox/lab.com.pem.crt
```

A successful result ends with `OK`. Use `-verify_ip 127.0.0.1` instead of `-verify_hostname lab.com` to check the IP address.

Passing these checks confirms the certificate and configured trust relationship. You still need to configure trust in the actual client and ensure the server serves the intended certificate.

## Configuration

Every option supports an environment variable and a matching `--NAME=value` flag. On current source, precedence is:

**Non-empty explicit CLI value → environment variable → built-in default.**

An omitted or empty CLI value falls back to the environment/default. Pass each flag separately and quote values containing spaces or wildcard characters:

```bash
CERT_DNS="lab.com,127.0.0.1" ./certs-maker \
  --EXPIRE_DAYS=30 "--CERT_CN=Lab Server" --CUSTOM_FILE_NAME=lab-server
```

### Common options

| Option | Default | Purpose |
| --- | --- | --- |
| `CERT_DNS` | `lab.com,*.lab.com,*.data.lab.com` | Comma-separated DNS names, wildcard domains, and IP addresses. |
| `EXPIRE_DAYS` | `3650` | Certificate validity in days. Current source validates a positive integer. |
| `FOR_K8S` | `off` | Enable Kubernetes certificate extensions and naming. |
| `FOR_FIREFOX` | `off` | Generate a local root CA and signed server certificate. |
| `DIR` | `./ssl` | Output directory in the process’s filesystem. |
| `CUSTOM_FILE_NAME` | empty | Replace the output file prefix, without an extension. |

Use `ON` and `OFF` for mode flags. `on`, `true`, and `1` are true, ignoring case; other values are false.

`CERT_DNS` must contain at least one valid address. Invalid entries are discarded, duplicates are removed, and DNS names are lowercased. Bare `localhost` is not accepted by the normal input parser; Kubernetes mode adds it separately. `CERT_CN` changes the subject common name, not the SAN list.

Subject fields are configurable through `CERT_C`, `CERT_ST`, `CERT_L`, `CERT_O`, `CERT_OU`, and `CERT_CN`. The [complete configuration reference](docs/configuration.md) lists all defaults, expiry validation limits, input rules, and ownership options.

### Output and permissions

The default prefix is derived from the first valid address. Kubernetes mode adds `.k8s`; `CUSTOM_FILE_NAME` replaces the complete prefix. Use simple file names and relative directories such as `ssl` or `certs/server`: custom `DIR` values are normalized, so absolute paths, spaces, and case-sensitive names may change.

With Docker, `DIR` refers to a container path. Change the bind-mount source to choose another host location.

The container ownership helper requires `USER`, `UID`, and `GID` together. It uses Alpine account tools and makes output files readable by all users, including private keys. Review [ownership and key permissions](docs/configuration.md#file-ownership-and-private-keys) before enabling it. Manage ownership with your operating system’s tools on native installations.

## Replacing certificates

Each run generates new keys and overwrites files with the same output prefix. Save existing files elsewhere or use another directory/prefix before creating a replacement.

Firefox mode also regenerates the root CA. After replacement, deploy the new server certificate/key pair and update clients to trust the new public root certificate. An older root CA cannot verify certificates signed by the new one.

## Troubleshooting

- **Hostname or IP mismatch:** inspect SANs and include the address clients actually use in `CERT_DNS`. A common name alone does not establish hostname coverage.
- **Unknown issuer / untrusted certificate:** configure the trust store used by that client. For Firefox mode, import the root CA from the same generation as the deployed server certificate.
- **Custom validity ignored or Firefox verification fails on v3.8.1:** use current source or a newer release containing the fixes.
- **Missing output:** check the bind mount, output directory, and command exit status. Generation errors can leave partial files; confirm the run succeeded before deploying them.

See the [certificate guide](docs/certificates.md) for certificate/key checks and additional deployment details.

## More information

- [Installation and release compatibility](docs/installation.md)
- [Certificate modes, verification, and trust](docs/certificates.md)
- [Complete configuration reference](docs/configuration.md)
- [Docker Hub](https://hub.docker.com/r/soulteary/certs-maker)
- [Creating a lightweight Docker tool image for self-signed certificates](https://soulteary.com/2022/10/22/make-docker-tools-image-with-only-3md-self-signed-certificate-certs-maker.html)
- [How to create and use self-signed certificates](https://soulteary.com/2021/02/06/how-to-make-and-use-a-self-signed-certificate.html)

## License

[MIT](LICENSE).
