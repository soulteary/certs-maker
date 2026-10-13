# Installation

## Docker

The published image includes the certificate generator, OpenSSL, and a shell. The v3.8.1 manifest supports Linux `amd64`, `arm64`, `arm/v6`, and `arm/v7`; Docker selects the matching platform.

```bash
docker run --rm -v "$(pwd)/ssl:/ssl" \
  -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" \
  soulteary/certs-maker:v3.8.1
```

The mount maps host `./ssl` to container `/ssl`. The container starts in `/`, so the default `DIR=./ssl` resolves to `/ssl`. For another host location, change the source side of the mount; keep the destination `/ssl` unless you also change the in-container output path.

## Native releases

Download an archive and the matching checksums file from [GitHub Releases](https://github.com/soulteary/certs-maker/releases). Native binaries still launch `openssl` through `sh`: both commands must be available in `PATH`. Building without CGO does not remove those runtime dependencies.

Choose the archive for your operating system and CPU:

| Platform | Archive suffix |
| --- | --- |
| macOS, Apple Silicon | `darwin_arm64.tar.gz` |
| macOS, Intel | `darwin_amd64.tar.gz` |
| Linux, x86-64 | `linux_amd64.tar.gz` |
| Linux, ARM64 | `linux_arm64.tar.gz` |
| Linux, 32-bit x86 | `linux_386.tar.gz` |
| Linux, ARMv6 / ARMv7 | `linux_armv6.tar.gz` / `linux_armv7.tar.gz` |

For example, on an Apple Silicon Mac:

```bash
curl -fL -o certs-maker_3.8.1_darwin_arm64.tar.gz \
  https://github.com/soulteary/certs-maker/releases/download/v3.8.1/certs-maker_3.8.1_darwin_arm64.tar.gz
curl -fL -o certs-maker_3.8.1_checksums.txt \
  https://github.com/soulteary/certs-maker/releases/download/v3.8.1/certs-maker_3.8.1_checksums.txt
awk '$2 == "certs-maker_3.8.1_darwin_arm64.tar.gz" { print }' \
  certs-maker_3.8.1_checksums.txt | shasum -a 256 -c - &&
tar -xzf certs-maker_3.8.1_darwin_arm64.tar.gz
openssl version
command -v sh
./certs-maker "--CERT_DNS=lab.com,127.0.0.1" --DIR=ssl
```

Check that checksum verification reports `OK` before extracting and running the archive. On Linux, select the matching archive and use `sha256sum -c -` instead of `shasum -a 256 -c -`.

The ownership helper uses Alpine-style `addgroup` and `adduser` commands. On native installations, manage file ownership with the operating system's own tools instead of setting `USER`, `UID`, and `GID`.

## Build from source

Install the Go version required by [`go.mod`](../go.mod) (currently Go 1.27.2), along with OpenSSL and a shell:

```bash
git clone https://github.com/soulteary/certs-maker.git
cd certs-maker
go build -o certs-maker .
./certs-maker "--CERT_DNS=lab.com,127.0.0.1" --DIR=ssl
```

Use `./certs-maker --help` to inspect available flags.

To package current source into a local Docker image using the same runtime Dockerfile as releases, build a binary for the chosen Linux platform first. This example targets Linux `amd64`:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o certs-maker .
docker build --platform=linux/amd64 -f docker/Dockerfile.goreleaser \
  -t certs-maker:local .
docker run --rm --platform=linux/amd64 -v "$(pwd)/ssl:/ssl" \
  -e "CERT_DNS=lab.com,*.lab.com,*.data.lab.com" \
  certs-maker:local
```

For an ARM64 host, change `GOARCH` to `arm64` and both Docker platform arguments to `linux/arm64`.

## Release compatibility

The README pins v3.8.1 for reproducible basic examples. That published release predates the source fixes for:

- `EXPIRE_DAYS`: v3.8.1 signs certificates for 3650 days even when another value is provided.
- Firefox mode: v3.8.1 does not carry SAN extensions into the server certificate.
- Option precedence: in v3.8.1, an explicit CLI value equal to a built-in default may fail to override an environment variable; this includes `OFF` overriding `ON`.

Use a source build containing the relevant fix, or a later release that includes it, when you need these behaviors. A newer source checkout does not change an already published Docker tag. The [configuration reference](configuration.md) and [certificate guide](certificates.md) describe the corrected source behavior.

## Sizes

Release download sizes depend on the version, platform, and archive compression. Extracted binary sizes also depend on build flags. Docker registry transfer sizes are compressed layer sizes; local image sizes include unpacked layers. Check the selected release or image rather than relying on historical screenshots or a single fixed size.
