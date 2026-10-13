# Certs Maker

[![Go Test Coverage](./.github/coverage.svg)](./.github/go-test-report.md) [![Go Report Card](./.github/goreportcard.svg)](./.github/goreportcard-report.md) [![CodeQL](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/codeql.yml) [![Release](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml/badge.svg)](https://github.com/soulteary/certs-maker/actions/workflows/release.yaml) [![Docker Image](https://img.shields.io/docker/pulls/soulteary/certs-maker.svg)](https://hub.docker.com/r/soulteary/certs-maker)

[English](README.md) | [简体中文](README_CN.md)

<img src="logo.png" alt="Certs Maker logo" width="160">

A small certificate generator for development and self-hosted services. Run it with Docker, or use a native macOS or Linux binary. It supports multiple DNS names, wildcard domains, and IP addresses, with PEM and DER output.

Certificates need explicit client trust. For private production environments, manage trust, private keys, deployment, and renewal according to your own PKI policy. The tool generates files; it does not install certificates or renew them automatically.

## Quick start

With Docker installed and running:

```bash
docker run --rm -v "$(pwd)/ssl:/ssl" \
  -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" \
  soulteary/certs-maker:v3.8.1
```

This writes the following files to `./ssl` on the host:

```text
ssl/
├── lab.com.conf
├── lab.com.der.crt
├── lab.com.der.key
├── lab.com.pem.crt
└── lab.com.pem.key
```

Use `lab.com.pem.crt` and `lab.com.pem.key` as the certificate/key pair for a server that accepts PEM files. DER files contain the same certificate and key in binary form. Keep private keys private.

Inspect the generated certificate with OpenSSL:

```bash
openssl x509 -in ssl/lab.com.pem.crt -noout -subject -issuer -dates -ext subjectAltName
```

Generating the files does not make a browser trust them. See the [certificate usage and trust guide](docs/certificates.md) for verification, client trust, mode-specific files, and Kubernetes use.

The commands pin the published v3.8.1 image. See [release compatibility](docs/installation.md#release-compatibility) before using custom expiry or Firefox mode; those fixes require a source build until a newer release includes them.

## Docker Compose

Save this as `docker-compose.yml`, or use the file in this repository:

```yaml
services:
  certs-maker:
    image: soulteary/certs-maker:v3.8.1
    environment:
      CERT_DNS: "lab.com,*.lab.com,*.data.lab.com"
    volumes:
      - ./ssl:/ssl
```

Run it once and remove the container afterward:

```bash
docker compose run --rm certs-maker
```

Legacy Compose installations can use `docker-compose run --rm certs-maker`.

## Certificate modes

- **Default:** generates a self-signed certificate with `CA:TRUE`, the requested SANs, and PEM/DER files.
- **Kubernetes (`FOR_K8S=ON`):** uses server certificate extensions with `CA:FALSE`, adds `*` and `localhost` to the SAN list, and adds `.k8s` to the default file prefix. You still create and deploy the Kubernetes Secret yourself.
- **Firefox (`FOR_FIREFOX=ON`):** generates a local root CA and a separately signed server certificate. Import the public root CA certificate into the clients that should trust it.

If both mode flags are enabled, Kubernetes mode takes precedence. [Mode details, output files, and verification](docs/certificates.md#certificate-modes) explain which certificate to deploy and which certificate to trust.

## Configuration

Every option supports an environment variable and a matching CLI flag. Pass each flag as a separate argument, and quote values containing spaces or wildcard characters:

```bash
./certs-maker "--CERT_DNS=lab.com,127.0.0.1" "--CERT_CN=Lab Server" --FOR_K8S=ON
```

See the [configuration reference](docs/configuration.md) for defaults, accepted values, precedence, file naming, output path limitations, and ownership settings.

## Native installation and source builds

[GitHub Releases](https://github.com/soulteary/certs-maker/releases) provides macOS and Linux binaries plus SHA-256 checksums. Native execution needs `openssl` and `sh` in `PATH`; the Docker image includes these dependencies.

The [installation guide](docs/installation.md) covers platform selection, checksum verification, source builds, and local Docker builds. File and image sizes vary by version and architecture; compressed download size differs from extracted binary or local image size.

## Common questions

- **The browser still reports an untrusted certificate.** Configure client trust and verify the hostname; see [trust and verification](docs/certificates.md#client-trust).
- **Where are my files?** The bind mount maps host `./ssl` to container `/ssl`. `DIR` changes the path inside the container; see [output paths](docs/configuration.md#output-paths-and-file-names).
- **Can I rerun the command?** Yes, but files with the same prefix are overwritten and new keys are generated. Move existing files aside or use another prefix first.
- **Why is custom expiry ignored or Firefox verification failing on v3.8.1?** That release predates the fixes; use a [source build](docs/installation.md#build-from-source).

## Related documentation

- [Creating a lightweight Docker tool image for self-signed certificates](https://soulteary.com/2022/10/22/make-docker-tools-image-with-only-3md-self-signed-certificate-certs-maker.html)
- [How to create and use self-signed certificates](https://soulteary.com/2021/02/06/how-to-make-and-use-a-self-signed-certificate.html)
- [Docker Hub](https://hub.docker.com/r/soulteary/certs-maker)
- [Security policy](SECURITY.md) · [MIT license](LICENSE)
